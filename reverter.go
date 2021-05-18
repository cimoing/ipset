package ipset

import (
	"net"
	"strings"

	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("ipset")

// ResponseReverter reverses the operations done on the question section of a packet.
// This is need because the client will otherwise disregards the response, i.e.
// dig will complain with ';; Question section mismatch: got example.org/HINFO/IN'
type ResponseReverter struct {
	dns.ResponseWriter
	originalQuestion dns.Question
	ipset            *IPSet
}

// NewResponseReverter returns a pointer to a new ResponseReverter.
func NewResponseReverter(w dns.ResponseWriter, r *dns.Msg, ipset *IPSet) *ResponseReverter {
	return &ResponseReverter{
		ResponseWriter:   w,
		originalQuestion: r.Question[0],
		ipset:            ipset,
	}
}

// WriteMsg records the status code and calls the underlying ResponseWriter's WriteMsg method.
func (r *ResponseReverter) WriteMsg(res *dns.Msg) error {
	res.Question[0] = r.originalQuestion
	log.Info("Start Match ", r.originalQuestion.Name)
	if r.ipset.domains.Get(r.originalQuestion.Name[0:len(r.originalQuestion.Name)-1]) != 1 {
		log.Info("Match Failed", r.originalQuestion.Name)
		return r.ResponseWriter.WriteMsg(res)
	}
	log.Info("Matched ", r.originalQuestion.Name)
	for _, rr := range res.Answer {
		if rr.Header().Rrtype != dns.TypeA && rr.Header().Rrtype != dns.TypeAAAA {
			continue
		}

		ss := strings.Split(rr.String(), "\t")
		if len(ss) != 5 {
			continue
		}
		ip := net.ParseIP(ss[4])
		for _, listName := range r.ipset.listName {
			if err := addIP(ip, listName); err != nil {
				log.Error("add IP:", ip, " to ipset:", listName, ", result:", err)
			}
		}
	}
	return r.ResponseWriter.WriteMsg(res)
}

// Write is a wrapper that records the size of the message that gets written.
func (r *ResponseReverter) Write(buf []byte) (int, error) {
	n, err := r.ResponseWriter.Write(buf)
	return n, err
}
