# ipset
ipset plugin for CoreDNS

# syntax

包含``.baidu.com``,``.126.com``并排除``.ad.baidu.com``

<code>
ipset ipset_list {
    include .baidu.com .126.com
    exclude .ad.baidu.com
}
</code>
