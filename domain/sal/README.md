acl 防止rpc请求外部服务的防腐层

通过RPC集成

其中，被集成方（A上下文，U 是 Upstream 的缩写）采用了开放主机+发布语言的方式，而集成方（B上下文，D 是 Downstream 的缩写）则使用了防腐层的方式。

    OHS(Open Host Service)：开放主机服务，即定义一种协议，子系统可以通过该协议来访问你的服务。
    PL(Published Language)：发布语言，通常跟 OHS 一起使用，用于定义开放主机的协议。
    ACL(Anticorruption Layer)：防腐层，一个上下文通过一些适配和转换来跟另一上下文交互。

