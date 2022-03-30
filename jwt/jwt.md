nn mJWT认证

Json web token （JWT）基于json 的开放标准，适用于分布式站点的单点登录（SSO）场景



http 是无状态的，cookie 是存储在客户端，是一小块字符串

cookie 是不可以跨域额，会绑定cookie ，一级域名和二级域名是允许共享使用的 靠的是domain



cookie 重要性，  

cookie的重要性， name=value 键值对设置 名称 对应值。必须是字符串

name=value



domain  指定 cookie所属域名

 maxAge cookie失效的时间，单位秒。



secure 当secure为true时，安全协议是https。ssl等，默认为false



httpOnly 如果设置httpoly 属性，无法通过js 脚本，





## JWT

- jwt 服务器认证之后，生成一个json对象，发回给用户，

  ```
  {
  	"name": "张三",
  	"角色":"管理员",
  	"到期时间":"202了。。。"
  }
  ```

  这样sever 就不用保存任何session数据了，sever变成无状态的

- 现实中的 jwt 是一个很大的字符串，中间两个.分成三个部分

- Header 头部

- Payload 负载

- Signature 签名

Header.Payload.Signature



Header 一般是Json 对象，描述JWT 元数据 

```
{
	“alg”: "HS256",
	"typ": "JWT"
}
```

alg记录加密算法，type 表示 这个token 的type ，JWT令牌 统一写为 JWT

然后用base64URL 算法 转成字符串





Payload 部分是一个Json 对象，用来存放需要传递的数据，JWT 规定了7个官方子u

```
iss (issuer): 签发人
exp (EXPIRATION time): 过期时间
sub (SUBJECT) 主题
aud (audientce)受众
nbf （not before） 生效时间
iat （issued at） 签发时间
jti (jwt ID): 编号
```

这一段默认不加密的，任何人可以读到，所以不要把密码放在这个部分，也要base64URL 转化字符串



signature 部分可以是对

服务端有密钥 （sercret，）然后指定签名算法默认（HMAC SHA256）.产生签名



拼起来就是jwt



- JWT的最后实现方式
- 可以用cookie 但是不能跨域
- 放在http请求头 信息 Authorization 