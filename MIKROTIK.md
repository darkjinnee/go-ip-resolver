```bash
/usr/local/bin/go-ip-resolver --config=/root/.config/ip-resolver

$ sudo curl -L "https://github.com/darkjinnee/go-ip-resolver/releases/download/v1.2.0/ip-resolver-linux" -o /usr/local/bin/go-ip-resolver
$ sudo chmod +x /usr/local/bin/go-ip-resolver
```

```bash
:local listName "l2tp-ips"
:local apiUrl "http://37.220.82.157:8080/resolve-flat-all-ipv4"
:local expireTime 3600
:local tmpFile "iplist.txt"

# удалить старый файл, если есть
:if ([:len [/file find name=$tmpFile]] > 0) do={ /file remove $tmpFile }

# скачать список IP
/tool fetch url=$apiUrl mode=http dst-path=$tmpFile

:delay 1
:local content [/file get $tmpFile contents]

# построчно обработать файл
:while ([:len $content] > 0) do={
:local nlPos [:find $content "\n"]
:if ($nlPos = nil) do={ :set nlPos [:len $content] }

    :local ip [:pick $content 0 $nlPos]
    :set content [:pick $content ($nlPos + 1) [:len $content]]

    :if ([:len $ip] > 0) do={
        :local findIP [/ip firewall address-list find list=$listName address=$ip]
        :if ($findIP = "") do={
            /ip firewall address-list add list=$listName address=$ip timeout=$expireTime
            :log info ("VPN-IPs: added " . $ip)
        } else={
            /ip firewall address-list set $findIP timeout=$expireTime
            :log info ("VPN-IPs: refreshed " . $ip)
        }
    }
}

/file remove $tmpFile
```