### 文件目录

* api ：提供对外通信接口
* cmd ：可执行脚本
* config ：配置相关
* encoding ：各类信息解析（json，proto，xml，yaml
* errors ：错误信息配置，记录
* internal ：项目相关功能内容
* log ：日志相关

### 相关目录详解

#### config
通过加载文件的方式，将需要的配置文件解析成 map[string]interface()，通过第三包[imdario/mergo](github.com/imdario/mergo) 将读取的文件信息赋值，通过 automic 进行原子性赋值加入到 sync.Map 当中方便后续读取；



    