-- MySQL dump 10.13  Distrib 5.7.44, for Linux (x86_64)
--
-- Host: localhost    Database: blogx
-- ------------------------------------------------------
-- Server version	5.7.44-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `article_digg_models`
--
use blogx;
DROP TABLE IF EXISTS `article_digg_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `article_digg_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `article_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`user_id`,`article_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `article_digg_models`
--

LOCK TABLES `article_digg_models` WRITE;
/*!40000 ALTER TABLE `article_digg_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `article_digg_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `article_models`
--

DROP TABLE IF EXISTS `article_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `article_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `title` varchar(32) DEFAULT NULL,
  `abstract` varchar(256) DEFAULT NULL,
  `content` longtext,
  `category_id` bigint(20) unsigned DEFAULT NULL,
  `tag_list` longtext,
  `cover` varchar(256) DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `look_count` bigint(20) DEFAULT NULL,
  `digg_count` bigint(20) DEFAULT NULL,
  `comment_count` bigint(20) DEFAULT NULL,
  `collect_count` bigint(20) DEFAULT NULL,
  `open_comment` tinyint(1) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_article_models_category_id` (`category_id`),
  KEY `idx_article_models_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `article_models`
--

LOCK TABLES `article_models` WRITE;
/*!40000 ALTER TABLE `article_models` DISABLE KEYS */;
INSERT INTO `article_models` VALUES (1,'2024-12-17 20:31:55.947','2024-12-17 20:31:55.947','1.go环境搭建','相对于其他语言，go语言的环境搭建及其简单\n\n官网\n\nhttps://go.dev/dl/\n\n访问不了的就访问中文网就好了\n\ngo安装包下载\n\nhttps://studygolang.com/dl\n\n','相对于其他语言，go语言的环境搭建及其简单\n\n\n\n\n\n官网\n\nhttps://go.dev/dl/\n\n\n\n访问不了的就访问中文网就好了\n\n\n\ngo安装包下载\n\nhttps://studygolang.com/dl\n\n\n\n安装指定版本的安装包就好了\n\n## windows下的安装\n\nwindows就选 `windows-arm64.zip`就好了\n\n![](https://image.fengfengzhidao.com/pic/20231105132441.png)\n\n然后需要将go的对应bin目录设置为环境变量，这一步是方便可以在命令行里面直接使用go命令\n\n还需要将go的第三方bin目录设置为环境变量，一般是在用户目录下，这一步是为了以后使用go install安装的第三方可执行文件可以直接使用\n\n\n\n![](https://image.fengfengzhidao.com/pic/20231105132405.png)\n\n\n\n\n\n## linux下的安装\n\n![](https://image.fengfengzhidao.com/pic/20231105132836.png)\n\n\n\n```Go\ncd /opt\nwget https://studygolang.com/dl/golang/go1.21.3.linux-amd64.tar.gz\ntar -xvf go1.21.3.linux-amd64.tar.gz\n\n```\n\n编辑环境变量\n\n```Vim script\nvim /etc/profile\n\n在文件后追加以下内容\nexport GOPROXY=https://goproxy.cn\nexport GOROOT=/opt/go\nexport PATH=$PATH:$GOROOT/bin\nexport GOPATH=/opt/go/pkg\nexport PATH=$PATH:$GOPATH/bin\n\n退出并保存，刷新环境变量\nsource /etc/profile\n\n```\n\n\n\n\n\n## 开发工具的选择\n\n理论上来说，用记事本也不是不行\n\n但是有一个趁手的兵器肯定还是更合适的\n\n首选肯定是goland，当然vscode也是可以的\n\nhttps://www.jetbrains.com/zh-cn/go/download/other.html\n\n\n\n2023.2版本之后的UI变成了类似vscode的风格，如果不喜欢的话，可以选2023.1之前的版本\n\n\n\n\n\nvscode下载\n\nhttps://code.visualstudio.com/\n\n\n\n然后去下载go的插件就好了\n\n\n\n\n\n\n\n## 参考文档\n\nvscode [https://code.visualstudio.com/](https://code.visualstudio.com/)\n\nvscode安装go环境 [https://blog.csdn.net/flurry_rain/article/details/128124573](https://blog.csdn.net/flurry_rain/article/details/128124573)\n\ngoland [https://www.jetbrains.com/zh-cn/go/download/other.html](https://www.jetbrains.com/zh-cn/go/download/other.html)\n\ngo安装 [https://studygolang.com/dl](https://studygolang.com/dl)\n\n\n\n\n\n\n\n',NULL,'go','/uploads/images001/de8f4a8a6336fb5742b382237dd988a4.png',1,0,0,0,0,1,3),(2,'2024-12-17 20:35:56.740','2024-12-17 20:35:56.740','2.变量与输入输出','变量定义\n\n标准的变量定义\n\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  // 先定义，再赋值\n  var name string\n  name = \"枫枫','## 变量定义\n\n标准的变量定义\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  // 先定义，再赋值\n  var name string\n  name = \"枫枫\"\n  fmt.Println(name)\n  \n  // var 变量名 类型 = \"变量值\"\n  var userName string = \"枫枫\"\n  fmt.Println(userName)\n}\n\n\n```\n\n> 如果一个变量定义了，但是没有赋值，那么这个变量的值就是这个类型的 `\"零值\"`\n\n\n\n### 变量类型省略\n\n```Go\nvar name = \"枫枫\"\n```\n\n### 简短声明\n\n```Go\nname := \"枫枫\"\n```\n\n\n\n\n\n### 全局变量与局部变量\n\n1. 定义在函数体（包括main函数）内的变量都是局部变量，定义了就必须使用\n2. 定义在外部的变量就是全局变量，可以只定义不使用\n\n```Go\npackage main\n\nimport \"fmt\"\n\nvar userName = \"枫枫知道\" // 可以不使用\n\nfunc main() {\n  // var 变量名 类型 = \"变量值\"\n  var name = \"枫枫\"\n  // 在函数体内定义的变量，必须要使用\n  fmt.Println(name)\n}\n\n```\n\n\n\n### 定义多个变量\n\n```Go\npackage main\n\nfunc main() {\n  var name1, name2, name3 string // 定义多个变量\n\n  var a1, a2 = \"枫枫\", \"知道\" // 定义多个变量并赋值\n  \n  a3, a4 := \"枫枫\", \"知道\" // 简短定义多个变量并赋值\n}\n\n```\n\n\n\n```Go\npackage main\n\nimport \"fmt\"\n\nvar (\n  name     string = \"枫枫\"\n  userName        = \"枫枫知道\"\n)\n\nfunc main() {\n  fmt.Println(name, userName)\n}\n\n\n```\n\n\n\n### 常量定义\n\n1. 定义的时候就要赋值\n2. 赋值之后就不能再修改了\n\n```Go\nconst name string = \"枫枫\" // 定义就要赋值\n//name = \"知道\" // 不能再修改了\nfmt.Println(name)\n```\n\n\n\n### 命名规范\n\n> 核心思想：首字母大写的变量、函数。方法，属性可在包外进行访问\n\n\n\n## 输出\n\n常用的输出函数\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  fmt.Println(\"枫枫知道\")\n  fmt.Println(1)\n  fmt.Println(true)\n  fmt.Println(\"什么\", \"都\", \"可以\", \"输出\")\n}\n\n```\n\n\n\n### 格式化输出\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  fmt.Printf(\"%v\\n\", \"你好\")          // 可以作为任何值的占位符输出\n  fmt.Printf(\"%v %T\\n\", \"枫枫\", \"枫枫\") // 打印类型\n  fmt.Printf(\"%d\\n\", 3)             // 整数\n  fmt.Printf(\"%.2f\\n\", 1.25)        // 小数\n  fmt.Printf(\"%s\\n\", \"哈哈哈\")         // 字符串\n  fmt.Printf(\"%#v\\n\", \"\")           // 用go的语法格式输出，很适合打印空字符串\n}\n\n```\n\n\n\n还有一个用的比较多的就是将格式化之后的内容赋值给一个变量\n\n```Go\nname := fmt.Sprintf(\"%v\", \"你好\")\nfmt.Println(name)\n```\n\n\n\n## 输入\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  fmt.Println(\"输入您的名字：\")\n  var name string\n  fmt.Scan(&name)  // 这里记住，要在变量的前面加个&, 后面讲指针会提到\n  fmt.Println(\"你输入的名字是\", name)\n}\n\n```\n\n\n\n\n\n## 参考文档\n\n1. 变量定义 [https://segmentfault.com/a/1190000042910744](https://segmentfault.com/a/1190000042910744?sort=votes)\n2. 输入函数  [https://www.php.cn/faq/507954.html](https://www.php.cn/faq/507954.html)\n3. 格式化 [https://www.php.cn/faq/527489.html](https://www.php.cn/faq/527489.html)',NULL,'go','/uploads/images001/52c512a2c01a51644ac979b05947cafe.png',2,0,0,0,0,1,3),(3,'2024-12-17 20:36:17.354','2024-12-17 20:36:17.354','3.基本数据类型','go语言的基本数据类型有\n\n\n整数形\n浮点型\n复数\n布尔\n字符串\n\n\n整数型\n\ngo语言的整数类型，具体细分有很多\n\nvar n1 uint8 = 2 \nvar n2 uint16 = 2 \nvar ','go语言的基本数据类型有\n\n1. 整数形\n2. 浮点型\n3. 复数\n4. 布尔\n5. 字符串\n\n\n\n## 整数型\n\ngo语言的整数类型，具体细分有很多\n\n```Go\nvar n1 uint8 = 2 \nvar n2 uint16 = 2 \nvar n3 uint32 = 2 \nvar n4 uint64 = 2 \nvar n5 uint = 2 \nvar n6 int8 = 2 \nvar n7 int16 = 2 \nvar n8 int32 = 2 \nvar n9 int64 = 2 \nvar n10 int = 2 \n```\n\n大家只需要记住以下几点\n\n1. 默认的数字定义类型是int类型\n2. 带个u就是无符号，只能存正整数\n3. 后面的数字就是2进制的位数\n4. uint8还有一个别名 byte， 一个字节=8个bit位\n5. int类型的大小取决于所使用的平台\n\n\n\n例如uint8，那就是8个二进制位，都用来存储数据，那最小就是0，最大就是2的八次方-1=255\n\n那int8，因为要拿一位存符合，使用实际只有七位可用，所以最小的就是负2的七次方=-128，最大的就是2的七次方-1=127\n\n至于为什么要减一，其实很好理解，因为实际到最后一个数字的时候，已经向前进位了，例如一个小时是60分钟，但是分钟最大只有59\n\n\n\n第五点的测试\n\n我是64位操作系统，那么我会试一下int是不是就是int64的最大上限\n\n2的63次方-1=9223372036854775807\n\n```Go\nfmt.Printf(\"%.0f\\n\", math.Pow(2, 63))\nvar n1 int = 9223372036854775807\nfmt.Println(n1)\nvar n2 int = 9223372036854775808 // 看它报不报错\nfmt.Println(n2)\n```\n\n\n\n## 浮点型\n\nGo语言支持两种浮点型数：float32 和 float64 \n\n1. float32 的浮点数的最大范围约为3.4e38，可以使用常量定义：`math.MaxFloat32`\n2. float64 的浮点数的最大范围约为 1.8e308，可以使用一个常量定义：`math.MaxFloat64`\n\n\n\n如果没有显式声明，则默认是float64\n\n\n\n\n\n## 字符型\n\n> 注意哦，是字符，不是字符串\n\n比较重要的两个类型是byte（单字节字符）和rune（多字节字符）\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  var c1 = \'a\'\n  var c2 = 97\n  fmt.Println(c1) // 直接打印都是数字\n  fmt.Println(c2)\n\n  fmt.Printf(\"%c %c\\n\", c1, c2) // 以字符的格式打印\n\n  var r1 rune = \'中\'\n  fmt.Printf(\"%c\\n\", r1)\n}\n\n```\n\n\n\n1. 在 Go 中，字符的本质是一个整数，直接输出时，是该字符对应的 UTF-8 编码的码值\n2. 可以直接给某个变量赋一个数字，然后按格式化输出时 %c ，会输出该数字对应的 unicode 字符  \n3. 字符类型是可以进行运算的，相当于一个整数，因为它都对应有 Unicode 码。\n\n\n\n## 字符串类型\n\n和字符不一样的是，字符的赋值是单引号，字符串的赋值是双引号\n\n```Go\nvar s string = \"枫枫知道\"\nfmt.Println(s)\n```\n\n\n\n### 转义字符\n\n一些常用的转义字符\n\n```Go\nfmt.Println(\"枫枫\\t知道\")              // 制表符\nfmt.Println(\"枫枫\\n知道\")              // 回车\nfmt.Println(\"\\\"枫枫\\\"知道\")            // 双引号\nfmt.Println(\"枫枫\\r知道\")              // 回到行首\nfmt.Println(\"C:\\\\pprof\\\\main.exe\") // 反斜杠\n```\n\n\n\n### 多行字符串\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  var s = `今天\n天气\n真好\n`\n  fmt.Println(s)\n}\n\n```\n\n在``这个里面，再出现转义字符就会原样输出了\n\n\n\n\n\n\n\n## 布尔类型\n\n布尔型数据只有 true（真）和 false（假）两个值\n\n1. 布尔类型变量的默认值为false\n2. Go 语言中不允许将整型强制转换为布尔型\n3. 布尔型无法参与数值运算，也无法与其他类型进行转换\n\n\n\n\n\n## 零值问题\n\n如果我们给一个基本数据类型只声明不赋值\n\n那么这个变量的值就是对应类型的零值，例如int就是0，bool就是false，字符串就是\"\"\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  var a1 int\n  var a2 float32\n  var a3 string\n  var a4 bool\n\n  fmt.Printf(\"%#v\\n\", a1)\n  fmt.Printf(\"%#v\\n\", a2)\n  fmt.Printf(\"%#v\\n\", a3)\n  fmt.Printf(\"%#v\\n\", a4)\n}\n\n```\n\n\n\n## 参考文档\n\n1. go语言类型 [https://blog.csdn.net/weixin_44211968/article/details/121221309](https://blog.csdn.net/weixin_44211968/article/details/121221309)\n2. golang int长度 [https://www.04ip.com/post/954343.html](https://www.04ip.com/post/954343.html)\n3. go int默认长度 [https://www.php.cn/faq/538086.html](https://www.php.cn/faq/538086.html)',NULL,'','',2,0,0,0,0,1,3),(4,'2024-12-17 20:37:17.834','2024-12-17 20:37:17.834','4.数组、切片、map','数组\n\n数组（Array）是一种非常常见的数据类型，几乎所有的计算机编程语言中都会用到它\n\n\n数组里的元素必须全部为同一类型，要嘛全部是字符串，要嘛全部是整数\n声明数组时，必须指定其长度或者大小\n\n\n','## 数组\n\n数组（Array）是一种非常常见的数据类型，几乎所有的计算机编程语言中都会用到它\n\n1. 数组里的元素必须全部为同一类型，要嘛全部是字符串，要嘛全部是整数\n2. 声明数组时，必须指定其长度或者大小\n\n```Go\npackage main\n\nimport \"fmt\" \n\nfunc main() {\n  var array [3]int = [3]int{1, 2, 3}\n  fmt.Println(array)\n  var array1 = [3]int{1, 2, 3}\n  fmt.Println(array1)\n  var array2 = [...]int{1, 2, 3}\n  fmt.Println(array2)\n}\n\n```\n\n如果要修改某个值，只能根据索引去找然后替换\n\n```Go\nvar array1 = [3]int{1, 2, 3}\narray1[0] = 10 // 根据索引找到对应的元素位置，然后替换\nfmt.Println(array1)\n```\n\n\n\n### 索引\n\n索引这个知识点，如果刚学编程，那么就得仔细听了，这个知识点到后面我就默认你已经会了\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  // 以定义一个字符串数组 a b c d 为例\n  var sList = []string{\"a\", \"b\", \"c\", \"d\"}\n  // 正向索引就从0开始，从左往右\n  /*\n    a     b     c     d\n    0     1     2     3\n  */\n  // 取值就通过索引去取\n  fmt.Println(sList[0]) // 拿到a这个元素\n  // 当然，有的语言支持负向索引，go不支持\n  /*\n      a      b     c     d\n      -4    -3    -2    -1\n  */\n  // 如果在go语言里面，我想拿到倒数第二个元素，怎么办？\n  fmt.Println(sList[len(sList)-2]) // 拿到a这个元素\n}\n\n```\n\n\n\n## 切片\n\n很明显啊，go里面的数组，长度被限制死了，所以不经常用\n\n所以go出了一个数组plus，叫做slice（切片）\n\n切片（Slice）相较于数组更灵活，因为在声明切片后其长度是可变的\n\n\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  // 定义一个字符串切片\n  var list []string\n\n  list = append(list, \"枫枫\")\n  list = append(list, \"知道\")\n  fmt.Println(list)\n  fmt.Println(len(list)) // 切片长度\n  // 修改第二个元素\n  list[1] = \"不知道\"\n  fmt.Println(list)\n}\n\n```\n\n\n\n\n\n### make函数\n\n除了基本数据类型，其他数据类型如果只定义不赋值，那么实际的值就是nil\n\n```Go\n// 定义一个字符串切片\nvar list []string\nfmt.Println(list == nil) // true\n```\n\n\n\n那么我们可以通过make函数创建指定长度，指定容量的切片了\n\n```Go\nmake([]type, length, capacity)\n```\n\n\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  // 定义一个字符串切片\n  var list = make([]string, 0)\n  fmt.Println(list, len(list), cap(list))\n  fmt.Println(list == nil) // false\n\n  list1 := make([]int, 2, 2)\n  fmt.Println(list1, len(list1), cap(list1))\n}\n\n```\n\n\n\n### 为什么叫切片\n\n因为切片是数组切出来的\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  var list = [...]string{\"a\", \"b\", \"c\"}\n  slices := list[:] // 左一刀，右一刀  变成了切片\n  fmt.Println(slices)\n  fmt.Println(list[1:2]) // b\n}\n\n```\n\n\n\n### 切片排序\n\n```Go\nvar list = []int{4, 5, 3, 2, 7}\nfmt.Println(\"排序前:\", list)\nsort.Ints(list)\nfmt.Println(\"升序:\", list)\nsort.Sort(sort.Reverse(sort.IntSlice(list)))\nfmt.Println(\"降序:\", list)\n```\n\n\n\n## map\n\nGo语言中的map(映射、字典)是一种内置的数据结构，它是一个`无序`的key-value对的集合\n\nmap的key必须是基本数据类型，value可以是任意类型\n\n> 注意，map使用之前一定要初始化\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  // 声明\n  var m1 map[string]string\n  // 初始化1\n  m1 = make(map[string]string)\n  // 初始化2\n  m1 = map[string]string{}\n  // 设置值\n  m1[\"name\"] = \"枫枫\"\n  fmt.Println(m1)\n  // 取值\n  fmt.Println(m1[\"name\"])\n  // 删除值\n  delete(m1, \"name\")\n  fmt.Println(m1)\n}\n\n```\n\n\n\n也可以声明并赋值\n\n```Go\n// 声明并赋值\nvar m1 = map[string]string{}\nvar m2 = make(map[string]string)\n```\n\n\n\n### map取值\n\n如果只有一个参数接，那这个参数就是值，如果没有，这个值就是类型的零值\n\n如果两个参数接，那第二个参数就是布尔值，表示是否有这个元素\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  // 声明并赋值\n  var m1 = map[string]int{\n    \"age\": 21,\n  }\n  age1 := m1[\"age1\"] // 取一个不存在的\n  fmt.Println(age1)\n  age2, ok := m1[\"age1\"]\n  fmt.Println(age2, ok)\n}\n\n```\n\n\n\n\n\n## 参考文档\n\n1. 数组与切片 [https://zhuanlan.zhihu.com/p/456844080](https://zhuanlan.zhihu.com/p/456844080)\n2. map [https://blog.csdn.net/qq_39312146/article/details/127914811](https://blog.csdn.net/qq_39312146/article/details/127914811)',NULL,'','/uploads/images001/1e6fdd4abee43f87320a2bc1f094acc7.png',2,0,0,0,0,1,3),(5,'2024-12-17 20:37:34.883','2024-12-17 20:37:34.883','5.判断语句','if语句\n\n以年龄为例，输入的年龄在某一个区间，就输出对应的提示信息\n\n<=0     未出生\n1-18    未成年\n18-35   青年\n>=35    中年\n\n\n很明显，这是一个多选一的情况\n\n','## if语句\n\n以年龄为例，输入的年龄在某一个区间，就输出对应的提示信息\n\n```Go\n<=0     未出生\n1-18    未成年\n18-35   青年\n>=35    中年\n```\n\n很明显，这是一个多选一的情况\n\n我们有很多中方式来实现\n\n中断式\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  fmt.Println(\"请输入你的年龄：\")\n  var age int\n  fmt.Scan(&age)\n\n  if age <= 0 {\n    fmt.Println(\"未出生\")\n    return\n  }\n  if age <= 18 {\n    fmt.Println(\"未成年\")\n    return\n  }\n  if age <= 35 {\n    fmt.Println(\"青年\")\n    return\n  }\n  fmt.Println(\"中年\")\n\n}\n\n```\n\n它也有个好听的名字，叫卫语句\n\n\n\n嵌套式\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  fmt.Println(\"请输入你的年龄：\")\n  var age int\n  fmt.Scan(&age)\n\n  if age <= 18 {\n    if age <= 0 {\n      fmt.Println(\"未出生\")\n    } else {\n      fmt.Println(\"未成年\")\n    }\n  } else {\n    if age <= 35 {\n      fmt.Println(\"青年\")\n    } else {\n      fmt.Println(\"中年\")\n    }\n  }\n}\n\n```\n\n\n\n多条件式\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  fmt.Println(\"请输入你的年龄：\")\n  var age int\n  fmt.Scan(&age)\n\n  if age <= 0 {\n    fmt.Println(\"未出生\")\n  }\n  if age > 0 && age <= 18 {\n    fmt.Println(\"未成年\")\n  }\n  if age > 18 && age <= 35 {\n    fmt.Println(\"青年\")\n  }\n  if age > 35 {\n    fmt.Println(\"中年\")\n  }\n}\n\n```\n\n\n\n\n\n\n\n## switch语句\n\n还是上面那个案例，如果是用switch就很直观了\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  fmt.Println(\"请输入你的年龄：\")\n  var age int\n  fmt.Scan(&age)\n\n  switch {\n  case age <= 0:\n    fmt.Println(\"未出生\")\n  case age <= 18:\n    fmt.Println(\"未成年\")\n  case age <= 35:\n    fmt.Println(\"青年\")\n  default:\n    fmt.Println(\"中年\")\n  }\n}\n\n```\n\n除了这样的写法，还有枚举所有的可能值\n\n例如\n\n```Go\n1   星期一\n2   星期二\n3   星期三\n```\n\n\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  fmt.Println(\"请输入星期数字：\")\n  var week int\n  fmt.Scan(&week)\n\n  switch week {\n  case 1:\n    fmt.Println(\"周一\")\n  case 2:\n    fmt.Println(\"周二\")\n  case 3:\n    fmt.Println(\"周三\")\n  case 4:\n    fmt.Println(\"周四\")\n  case 5:\n    fmt.Println(\"周五\")\n  case 6, 7:\n    fmt.Println(\"周末\")\n  default:\n    fmt.Println(\"错误\")\n  }\n}\n\n\n```\n\n可以理解为case的值就是switch的枚举结果\n\n\n\n一般来说，go的switch的多选一，满足其中一个结果之后，就结束switch了\n\n例如\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  fmt.Println(\"请输入你的年龄：\")\n  var age int\n  fmt.Scan(&age)\n\n  switch {\n  case age <= 0:\n    fmt.Println(\"未出生\")\n  case age <= 18:\n    fmt.Println(\"未成年\")\n  case age <= 35:\n    fmt.Println(\"青年\")\n  default:\n    fmt.Println(\"中年\")\n  }\n}\n```\n\n我输入一个12，我希望它能输出满足的所有条件，例如我希望它输出，未成年，青年\n\n```Go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n  fmt.Println(\"请输入你的年龄：\")\n  var age int\n  fmt.Scan(&age)\n\n  switch {\n  case age <= 0:\n    fmt.Println(\"未出生\")\n    fallthrough\n  case age <= 18:\n    fmt.Println(\"未成年\")\n    fallthrough\n  case age <= 35:\n    fmt.Println(\"青年\")\n  default:\n    fmt.Println(\"中年\")\n  }\n}\n\n```\n\n\n\n\n\n## 参考文档\n\n卫语句  [https://blog.csdn.net/wangpaiblog/article/details/114909737](https://blog.csdn.net/wangpaiblog/article/details/114909737)\n\nswitch [https://blog.csdn.net/qq_43470538/article/details/130436381](https://blog.csdn.net/qq_43470538/article/details/130436381)\n\nswitch [https://www.php.cn/faq/540744.html](https://www.php.cn/faq/540744.html)\n\n',NULL,'','',2,0,0,0,0,1,3);
/*!40000 ALTER TABLE `article_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `banner_models`
--

DROP TABLE IF EXISTS `banner_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `banner_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `show` tinyint(1) DEFAULT NULL,
  `cover` longtext,
  `href` longtext,
  `type` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `banner_models`
--

LOCK TABLES `banner_models` WRITE;
/*!40000 ALTER TABLE `banner_models` DISABLE KEYS */;
INSERT INTO `banner_models` VALUES (1,'2024-12-17 20:39:01.034','2024-12-17 20:59:39.596',1,'/uploads/images001/1f778548ddc30677cb309352a7dd9ff7.png','https://www.fengfengzhidao.com/',1),(2,'2024-12-17 20:45:15.003','2024-12-17 21:01:24.882',1,'/uploads/images001/59b2e1febf35b7e1322d8b16216d7a03.png','https://www.bilibili.com/video/BV1zYB1YaEVL',2);
/*!40000 ALTER TABLE `banner_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `category_models`
--

DROP TABLE IF EXISTS `category_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `category_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `title` varchar(32) DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `category_models`
--

LOCK TABLES `category_models` WRITE;
/*!40000 ALTER TABLE `category_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `category_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `chat_models`
--

DROP TABLE IF EXISTS `chat_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `chat_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `send_user_id` bigint(20) unsigned DEFAULT NULL,
  `rev_user_id` bigint(20) unsigned DEFAULT NULL,
  `msg_type` tinyint(4) DEFAULT NULL,
  `msg` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chat_models`
--

LOCK TABLES `chat_models` WRITE;
/*!40000 ALTER TABLE `chat_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `chat_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `collect_models`
--

DROP TABLE IF EXISTS `collect_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `collect_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `title` varchar(32) DEFAULT NULL,
  `abstract` varchar(256) DEFAULT NULL,
  `cover` varchar(256) DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `is_default` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `collect_models`
--

LOCK TABLES `collect_models` WRITE;
/*!40000 ALTER TABLE `collect_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `collect_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `comment_digg_models`
--

DROP TABLE IF EXISTS `comment_digg_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `comment_digg_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `comment_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`user_id`,`comment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comment_digg_models`
--

LOCK TABLES `comment_digg_models` WRITE;
/*!40000 ALTER TABLE `comment_digg_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `comment_digg_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `comment_models`
--

DROP TABLE IF EXISTS `comment_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `comment_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `content` varchar(256) DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `article_id` bigint(20) unsigned DEFAULT NULL,
  `parent_id` bigint(20) unsigned DEFAULT NULL,
  `root_parent_id` bigint(20) unsigned DEFAULT NULL,
  `digg_count` bigint(20) DEFAULT NULL,
  `apply_count` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_comment_models_user_id` (`user_id`),
  KEY `idx_comment_models_article_id` (`article_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comment_models`
--

LOCK TABLES `comment_models` WRITE;
/*!40000 ALTER TABLE `comment_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `comment_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `global_notification_models`
--

DROP TABLE IF EXISTS `global_notification_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `global_notification_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `title` varchar(32) DEFAULT NULL,
  `icon` varchar(256) DEFAULT NULL,
  `content` varchar(64) DEFAULT NULL,
  `href` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `global_notification_models`
--

LOCK TABLES `global_notification_models` WRITE;
/*!40000 ALTER TABLE `global_notification_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `global_notification_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `image_models`
--

DROP TABLE IF EXISTS `image_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `image_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `filename` varchar(64) DEFAULT NULL,
  `path` varchar(256) DEFAULT NULL,
  `size` bigint(20) DEFAULT NULL,
  `hash` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `image_models`
--

LOCK TABLES `image_models` WRITE;
/*!40000 ALTER TABLE `image_models` DISABLE KEYS */;
INSERT INTO `image_models` VALUES (1,'2024-12-17 20:28:58.110','2024-12-17 20:28:58.110','1.png','uploads/images001/62cc6c0d334c14e239294c14cc36684a.png',107596,'62cc6c0d334c14e239294c14cc36684a'),(2,'2024-12-17 20:31:33.881','2024-12-17 20:31:33.881','51220e5b-056f-4399-925a-3eb4ca9299d7.png','uploads/images001/de8f4a8a6336fb5742b382237dd988a4.png',150834,'de8f4a8a6336fb5742b382237dd988a4'),(3,'2024-12-17 20:33:07.975','2024-12-17 20:33:07.975','th.png','uploads/images001/06c11c28d0d4488e7d7cb51e9506aedf.png',75051,'06c11c28d0d4488e7d7cb51e9506aedf'),(4,'2024-12-17 20:35:29.952','2024-12-17 20:35:29.952','20231006164712__wallhaven-zyxvqy.png','uploads/images001/52c512a2c01a51644ac979b05947cafe.png',39618,'52c512a2c01a51644ac979b05947cafe'),(5,'2024-12-17 20:37:15.541','2024-12-17 20:37:15.541','20231123092743__1123-i.png','uploads/images001/1e6fdd4abee43f87320a2bc1f094acc7.png',53064,'1e6fdd4abee43f87320a2bc1f094acc7'),(6,'2024-12-17 20:38:53.695','2024-12-17 20:38:53.695','20231006164712__wallhaven-zyxvqy.png','uploads/images001/5cf2c7170f85a31698f5ed5e9580d457.png',39076,'5cf2c7170f85a31698f5ed5e9580d457'),(7,'2024-12-17 20:40:04.719','2024-12-17 20:40:04.719','gvb-banner.png','uploads/images001/1577b579f9715a5419eec623b2b124e2.png',23892,'1577b579f9715a5419eec623b2b124e2'),(8,'2024-12-17 20:42:54.782','2024-12-17 20:42:54.782','gvb-banner.png','uploads/images001/40886dbe8b3b2eda5e40b4c96976f3dc.png',13567,'40886dbe8b3b2eda5e40b4c96976f3dc'),(9,'2024-12-17 20:44:23.523','2024-12-17 20:44:23.523','20231123092743__1123-i.png','uploads/images001/d9c44a98e17dfe218bacf4969fd29c4b.png',33745,'d9c44a98e17dfe218bacf4969fd29c4b'),(10,'2024-12-17 20:45:12.084','2024-12-17 20:45:12.084','gvb-banner.png','uploads/images001/a1d54d533d3ec7fa5d0d6fe33fd243e0.png',12996,'a1d54d533d3ec7fa5d0d6fe33fd243e0'),(11,'2024-12-17 20:49:03.310','2024-12-17 20:49:03.310','20231123092743__1123-i.png','uploads/images001/f02b41d7bb728da686df65f1ffdfeca5.png',52445,'f02b41d7bb728da686df65f1ffdfeca5'),(12,'2024-12-17 20:55:41.297','2024-12-17 20:55:41.297','xx1.png','uploads/images001/e9b2994d6de421c6d8e81f1c7acdfc26.png',12865,'e9b2994d6de421c6d8e81f1c7acdfc26'),(13,'2024-12-17 20:59:34.310','2024-12-17 20:59:34.310','xx1.png','uploads/images001/1f778548ddc30677cb309352a7dd9ff7.png',992618,'1f778548ddc30677cb309352a7dd9ff7'),(14,'2024-12-17 21:01:15.101','2024-12-17 21:01:15.101','blogx_server.png','uploads/images001/59b2e1febf35b7e1322d8b16216d7a03.png',597830,'59b2e1febf35b7e1322d8b16216d7a03');
/*!40000 ALTER TABLE `image_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_models`
--

DROP TABLE IF EXISTS `log_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `log_type` tinyint(4) DEFAULT NULL,
  `title` varchar(64) DEFAULT NULL,
  `content` longtext,
  `level` tinyint(4) DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `ip` varchar(32) DEFAULT NULL,
  `addr` varchar(64) DEFAULT NULL,
  `is_read` tinyint(1) DEFAULT NULL,
  `login_status` tinyint(1) DEFAULT NULL,
  `username` varchar(32) DEFAULT NULL,
  `pwd` varchar(32) DEFAULT NULL,
  `login_type` tinyint(4) DEFAULT NULL,
  `service_name` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_models`
--

LOCK TABLES `log_models` WRITE;
/*!40000 ALTER TABLE `log_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `message_models`
--

DROP TABLE IF EXISTS `message_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `message_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `type` tinyint(4) DEFAULT NULL,
  `rev_user_id` bigint(20) unsigned DEFAULT NULL,
  `action_user_id` bigint(20) unsigned DEFAULT NULL,
  `action_user_nickname` longtext,
  `action_user_avatar` longtext,
  `title` longtext,
  `content` longtext,
  `article_id` bigint(20) unsigned DEFAULT NULL,
  `article_title` longtext,
  `comment_id` bigint(20) unsigned DEFAULT NULL,
  `link_title` longtext,
  `link_href` longtext,
  `is_read` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `message_models`
--

LOCK TABLES `message_models` WRITE;
/*!40000 ALTER TABLE `message_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `message_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `site_flow_models`
--

DROP TABLE IF EXISTS `site_flow_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `site_flow_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `count` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `site_flow_models`
--

LOCK TABLES `site_flow_models` WRITE;
/*!40000 ALTER TABLE `site_flow_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `site_flow_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `text_models`
--

DROP TABLE IF EXISTS `text_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `text_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `article_id` bigint(20) unsigned DEFAULT NULL,
  `head` longtext,
  `body` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `text_models`
--

LOCK TABLES `text_models` WRITE;
/*!40000 ALTER TABLE `text_models` DISABLE KEYS */;
INSERT INTO `text_models` VALUES (1,'2024-12-17 20:31:55.948','2024-12-17 20:31:55.948',1,'1.go环境搭建','相对于其他语言，go语言的环境搭建及其简单官网https://go.dev/dl/访问不了的就访问中文网就好了go安装包下载https://studygolang.com/dl安装指定版本的安装包就好了'),(2,'2024-12-17 20:31:55.948','2024-12-17 20:31:55.948',1,'windows下的安装','windows就选 `windows-arm64.zip`就好了![](https://image.fengfengzhidao.com/pic/20231105132441.png)然后需要将go的对应bin目录设置为环境变量，这一步是方便可以在命令行里面直接使用go命令还需要将go的第三方bin目录设置为环境变量，一般是在用户目录下，这一步是为了以后使用go install安装的第三方可执行文件可以直接使用![](https://image.fengfengzhidao.com/pic/20231105132405.png)'),(3,'2024-12-17 20:31:55.948','2024-12-17 20:31:55.948',1,'linux下的安装','![](https://image.fengfengzhidao.com/pic/20231105132836.png)```Gocd /optwget https://studygolang.com/dl/golang/go1.21.3.linux-amd64.tar.gztar -xvf go1.21.3.linux-amd64.tar.gz```编辑环境变量```Vim scriptvim /etc/profile在文件后追加以下内容export GOPROXY=https://goproxy.cnexport GOROOT=/opt/goexport PATH=$PATH:$GOROOT/binexport GOPATH=/opt/go/pkgexport PATH=$PATH:$GOPATH/bin退出并保存，刷新环境变量source /etc/profile```'),(4,'2024-12-17 20:31:55.948','2024-12-17 20:31:55.948',1,'开发工具的选择','理论上来说，用记事本也不是不行但是有一个趁手的兵器肯定还是更合适的首选肯定是goland，当然vscode也是可以的https://www.jetbrains.com/zh-cn/go/download/other.html2023.2版本之后的UI变成了类似vscode的风格，如果不喜欢的话，可以选2023.1之前的版本vscode下载https://code.visualstudio.com/然后去下载go的插件就好了'),(5,'2024-12-17 20:31:55.948','2024-12-17 20:31:55.948',1,'参考文档','vscode [https://code.visualstudio.com/](https://code.visualstudio.com/)vscode安装go环境 [https://blog.csdn.net/flurry_rain/article/details/128124573](https://blog.csdn.net/flurry_rain/article/details/128124573)goland [https://www.jetbrains.com/zh-cn/go/download/other.html](https://www.jetbrains.com/zh-cn/go/download/other.html)go安装 [https://studygolang.com/dl](https://studygolang.com/dl)'),(6,'2024-12-17 20:35:56.742','2024-12-17 20:35:56.742',2,'2.变量与输入输出',''),(7,'2024-12-17 20:35:56.742','2024-12-17 20:35:56.742',2,'变量定义','标准的变量定义```Gopackage mainimport \"fmt\"func main() {  // 先定义，再赋值  var name string  name = \"枫枫\"  fmt.Println(name)    // var 变量名 类型 = \"变量值\"  var userName string = \"枫枫\"  fmt.Println(userName)}```> 如果一个变量定义了，但是没有赋值，那么这个变量的值就是这个类型的 `\"零值\"`'),(8,'2024-12-17 20:35:56.742','2024-12-17 20:35:56.742',2,'变量类型省略','```Govar name = \"枫枫\"```'),(9,'2024-12-17 20:35:56.742','2024-12-17 20:35:56.742',2,'简短声明','```Goname := \"枫枫\"```'),(10,'2024-12-17 20:35:56.742','2024-12-17 20:35:56.742',2,'全局变量与局部变量','1. 定义在函数体（包括main函数）内的变量都是局部变量，定义了就必须使用2. 定义在外部的变量就是全局变量，可以只定义不使用```Gopackage mainimport \"fmt\"var userName = \"枫枫知道\" // 可以不使用func main() {  // var 变量名 类型 = \"变量值\"  var name = \"枫枫\"  // 在函数体内定义的变量，必须要使用  fmt.Println(name)}```'),(11,'2024-12-17 20:35:56.742','2024-12-17 20:35:56.742',2,'定义多个变量','```Gopackage mainfunc main() {  var name1, name2, name3 string // 定义多个变量  var a1, a2 = \"枫枫\", \"知道\" // 定义多个变量并赋值    a3, a4 := \"枫枫\", \"知道\" // 简短定义多个变量并赋值}``````Gopackage mainimport \"fmt\"var (  name     string = \"枫枫\"  userName        = \"枫枫知道\")func main() {  fmt.Println(name, userName)}```'),(12,'2024-12-17 20:35:56.742','2024-12-17 20:35:56.742',2,'常量定义','1. 定义的时候就要赋值2. 赋值之后就不能再修改了```Goconst name string = \"枫枫\" // 定义就要赋值//name = \"知道\" // 不能再修改了fmt.Println(name)```'),(13,'2024-12-17 20:35:56.742','2024-12-17 20:35:56.742',2,'命名规范','> 核心思想：首字母大写的变量、函数。方法，属性可在包外进行访问'),(14,'2024-12-17 20:35:56.742','2024-12-17 20:35:56.742',2,'输出','常用的输出函数```Gopackage mainimport \"fmt\"func main() {  fmt.Println(\"枫枫知道\")  fmt.Println(1)  fmt.Println(true)  fmt.Println(\"什么\", \"都\", \"可以\", \"输出\")}```'),(15,'2024-12-17 20:35:56.742','2024-12-17 20:35:56.742',2,'格式化输出','```Gopackage mainimport \"fmt\"func main() {  fmt.Printf(\"%v\\n\", \"你好\")          // 可以作为任何值的占位符输出  fmt.Printf(\"%v %T\\n\", \"枫枫\", \"枫枫\") // 打印类型  fmt.Printf(\"%d\\n\", 3)             // 整数  fmt.Printf(\"%.2f\\n\", 1.25)        // 小数  fmt.Printf(\"%s\\n\", \"哈哈哈\")         // 字符串  fmt.Printf(\"%#v\\n\", \"\")           // 用go的语法格式输出，很适合打印空字符串}```还有一个用的比较多的就是将格式化之后的内容赋值给一个变量```Goname := fmt.Sprintf(\"%v\", \"你好\")fmt.Println(name)```'),(16,'2024-12-17 20:35:56.742','2024-12-17 20:35:56.742',2,'输入','```Gopackage mainimport \"fmt\"func main() {  fmt.Println(\"输入您的名字：\")  var name string  fmt.Scan(&name)  // 这里记住，要在变量的前面加个&, 后面讲指针会提到  fmt.Println(\"你输入的名字是\", name)}```'),(17,'2024-12-17 20:35:56.742','2024-12-17 20:35:56.742',2,'参考文档','1. 变量定义 [https://segmentfault.com/a/1190000042910744](https://segmentfault.com/a/1190000042910744?sort=votes)2. 输入函数  [https://www.php.cn/faq/507954.html](https://www.php.cn/faq/507954.html)3. 格式化 [https://www.php.cn/faq/527489.html](https://www.php.cn/faq/527489.html)'),(18,'2024-12-17 20:36:17.355','2024-12-17 20:36:17.355',3,'3.基本数据类型','go语言的基本数据类型有1. 整数形2. 浮点型3. 复数4. 布尔5. 字符串'),(19,'2024-12-17 20:36:17.355','2024-12-17 20:36:17.355',3,'整数型','go语言的整数类型，具体细分有很多```Govar n1 uint8 = 2 var n2 uint16 = 2 var n3 uint32 = 2 var n4 uint64 = 2 var n5 uint = 2 var n6 int8 = 2 var n7 int16 = 2 var n8 int32 = 2 var n9 int64 = 2 var n10 int = 2 ```大家只需要记住以下几点1. 默认的数字定义类型是int类型2. 带个u就是无符号，只能存正整数3. 后面的数字就是2进制的位数4. uint8还有一个别名 byte， 一个字节=8个bit位5. int类型的大小取决于所使用的平台例如uint8，那就是8个二进制位，都用来存储数据，那最小就是0，最大就是2的八次方-1=255那int8，因为要拿一位存符合，使用实际只有七位可用，所以最小的就是负2的七次方=-128，最大的就是2的七次方-1=127至于为什么要减一，其实很好理解，因为实际到最后一个数字的时候，已经向前进位了，例如一个小时是60分钟，但是分钟最大只有59第五点的测试我是64位操作系统，那么我会试一下int是不是就是int64的最大上限2的63次方-1=9223372036854775807```Gofmt.Printf(\"%.0f\\n\", math.Pow(2, 63))var n1 int = 9223372036854775807fmt.Println(n1)var n2 int = 9223372036854775808 // 看它报不报错fmt.Println(n2)```'),(20,'2024-12-17 20:36:17.355','2024-12-17 20:36:17.355',3,'浮点型','Go语言支持两种浮点型数：float32 和 float64 1. float32 的浮点数的最大范围约为3.4e38，可以使用常量定义：`math.MaxFloat32`2. float64 的浮点数的最大范围约为 1.8e308，可以使用一个常量定义：`math.MaxFloat64`如果没有显式声明，则默认是float64'),(21,'2024-12-17 20:36:17.355','2024-12-17 20:36:17.355',3,'字符型','> 注意哦，是字符，不是字符串比较重要的两个类型是byte（单字节字符）和rune（多字节字符）```Gopackage mainimport \"fmt\"func main() {  var c1 = \'a\'  var c2 = 97  fmt.Println(c1) // 直接打印都是数字  fmt.Println(c2)  fmt.Printf(\"%c %c\\n\", c1, c2) // 以字符的格式打印  var r1 rune = \'中\'  fmt.Printf(\"%c\\n\", r1)}```1. 在 Go 中，字符的本质是一个整数，直接输出时，是该字符对应的 UTF-8 编码的码值2. 可以直接给某个变量赋一个数字，然后按格式化输出时 %c ，会输出该数字对应的 unicode 字符  3. 字符类型是可以进行运算的，相当于一个整数，因为它都对应有 Unicode 码。'),(22,'2024-12-17 20:36:17.355','2024-12-17 20:36:17.355',3,'字符串类型','和字符不一样的是，字符的赋值是单引号，字符串的赋值是双引号```Govar s string = \"枫枫知道\"fmt.Println(s)```'),(23,'2024-12-17 20:36:17.355','2024-12-17 20:36:17.355',3,'转义字符','一些常用的转义字符```Gofmt.Println(\"枫枫\\t知道\")              // 制表符fmt.Println(\"枫枫\\n知道\")              // 回车fmt.Println(\"\\\"枫枫\\\"知道\")            // 双引号fmt.Println(\"枫枫\\r知道\")              // 回到行首fmt.Println(\"C:\\\\pprof\\\\main.exe\") // 反斜杠```'),(24,'2024-12-17 20:36:17.355','2024-12-17 20:36:17.355',3,'多行字符串','```Gopackage mainimport \"fmt\"func main() {  var s = `今天天气真好`  fmt.Println(s)}```在``这个里面，再出现转义字符就会原样输出了'),(25,'2024-12-17 20:36:17.355','2024-12-17 20:36:17.355',3,'布尔类型','布尔型数据只有 true（真）和 false（假）两个值1. 布尔类型变量的默认值为false2. Go 语言中不允许将整型强制转换为布尔型3. 布尔型无法参与数值运算，也无法与其他类型进行转换'),(26,'2024-12-17 20:36:17.355','2024-12-17 20:36:17.355',3,'零值问题','如果我们给一个基本数据类型只声明不赋值那么这个变量的值就是对应类型的零值，例如int就是0，bool就是false，字符串就是\"\"```Gopackage mainimport \"fmt\"func main() {  var a1 int  var a2 float32  var a3 string  var a4 bool  fmt.Printf(\"%#v\\n\", a1)  fmt.Printf(\"%#v\\n\", a2)  fmt.Printf(\"%#v\\n\", a3)  fmt.Printf(\"%#v\\n\", a4)}```'),(27,'2024-12-17 20:36:17.355','2024-12-17 20:36:17.355',3,'参考文档','1. go语言类型 [https://blog.csdn.net/weixin_44211968/article/details/121221309](https://blog.csdn.net/weixin_44211968/article/details/121221309)2. golang int长度 [https://www.04ip.com/post/954343.html](https://www.04ip.com/post/954343.html)3. go int默认长度 [https://www.php.cn/faq/538086.html](https://www.php.cn/faq/538086.html)'),(28,'2024-12-17 20:37:17.836','2024-12-17 20:37:17.836',4,'4.数组、切片、map',''),(29,'2024-12-17 20:37:17.836','2024-12-17 20:37:17.836',4,'数组','数组（Array）是一种非常常见的数据类型，几乎所有的计算机编程语言中都会用到它1. 数组里的元素必须全部为同一类型，要嘛全部是字符串，要嘛全部是整数2. 声明数组时，必须指定其长度或者大小```Gopackage mainimport \"fmt\" func main() {  var array [3]int = [3]int{1, 2, 3}  fmt.Println(array)  var array1 = [3]int{1, 2, 3}  fmt.Println(array1)  var array2 = [...]int{1, 2, 3}  fmt.Println(array2)}```如果要修改某个值，只能根据索引去找然后替换```Govar array1 = [3]int{1, 2, 3}array1[0] = 10 // 根据索引找到对应的元素位置，然后替换fmt.Println(array1)```'),(30,'2024-12-17 20:37:17.836','2024-12-17 20:37:17.836',4,'索引','索引这个知识点，如果刚学编程，那么就得仔细听了，这个知识点到后面我就默认你已经会了```Gopackage mainimport \"fmt\"func main() {  // 以定义一个字符串数组 a b c d 为例  var sList = []string{\"a\", \"b\", \"c\", \"d\"}  // 正向索引就从0开始，从左往右  /*    a     b     c     d    0     1     2     3  */  // 取值就通过索引去取  fmt.Println(sList[0]) // 拿到a这个元素  // 当然，有的语言支持负向索引，go不支持  /*      a      b     c     d      -4    -3    -2    -1  */  // 如果在go语言里面，我想拿到倒数第二个元素，怎么办？  fmt.Println(sList[len(sList)-2]) // 拿到a这个元素}```'),(31,'2024-12-17 20:37:17.836','2024-12-17 20:37:17.836',4,'切片','很明显啊，go里面的数组，长度被限制死了，所以不经常用所以go出了一个数组plus，叫做slice（切片）切片（Slice）相较于数组更灵活，因为在声明切片后其长度是可变的```Gopackage mainimport \"fmt\"func main() {  // 定义一个字符串切片  var list []string  list = append(list, \"枫枫\")  list = append(list, \"知道\")  fmt.Println(list)  fmt.Println(len(list)) // 切片长度  // 修改第二个元素  list[1] = \"不知道\"  fmt.Println(list)}```'),(32,'2024-12-17 20:37:17.836','2024-12-17 20:37:17.836',4,'make函数','除了基本数据类型，其他数据类型如果只定义不赋值，那么实际的值就是nil```Go// 定义一个字符串切片var list []stringfmt.Println(list == nil) // true```那么我们可以通过make函数创建指定长度，指定容量的切片了```Gomake([]type, length, capacity)``````Gopackage mainimport \"fmt\"func main() {  // 定义一个字符串切片  var list = make([]string, 0)  fmt.Println(list, len(list), cap(list))  fmt.Println(list == nil) // false  list1 := make([]int, 2, 2)  fmt.Println(list1, len(list1), cap(list1))}```'),(33,'2024-12-17 20:37:17.836','2024-12-17 20:37:17.836',4,'为什么叫切片','因为切片是数组切出来的```Gopackage mainimport \"fmt\"func main() {  var list = [...]string{\"a\", \"b\", \"c\"}  slices := list[:] // 左一刀，右一刀  变成了切片  fmt.Println(slices)  fmt.Println(list[1:2]) // b}```'),(34,'2024-12-17 20:37:17.836','2024-12-17 20:37:17.836',4,'切片排序','```Govar list = []int{4, 5, 3, 2, 7}fmt.Println(\"排序前:\", list)sort.Ints(list)fmt.Println(\"升序:\", list)sort.Sort(sort.Reverse(sort.IntSlice(list)))fmt.Println(\"降序:\", list)```'),(35,'2024-12-17 20:37:17.836','2024-12-17 20:37:17.836',4,'map','Go语言中的map(映射、字典)是一种内置的数据结构，它是一个`无序`的key-value对的集合map的key必须是基本数据类型，value可以是任意类型> 注意，map使用之前一定要初始化```Gopackage mainimport \"fmt\"func main() {  // 声明  var m1 map[string]string  // 初始化1  m1 = make(map[string]string)  // 初始化2  m1 = map[string]string{}  // 设置值  m1[\"name\"] = \"枫枫\"  fmt.Println(m1)  // 取值  fmt.Println(m1[\"name\"])  // 删除值  delete(m1, \"name\")  fmt.Println(m1)}```也可以声明并赋值```Go// 声明并赋值var m1 = map[string]string{}var m2 = make(map[string]string)```'),(36,'2024-12-17 20:37:17.836','2024-12-17 20:37:17.836',4,'map取值','如果只有一个参数接，那这个参数就是值，如果没有，这个值就是类型的零值如果两个参数接，那第二个参数就是布尔值，表示是否有这个元素```Gopackage mainimport \"fmt\"func main() {  // 声明并赋值  var m1 = map[string]int{    \"age\": 21,  }  age1 := m1[\"age1\"] // 取一个不存在的  fmt.Println(age1)  age2, ok := m1[\"age1\"]  fmt.Println(age2, ok)}```'),(37,'2024-12-17 20:37:17.836','2024-12-17 20:37:17.836',4,'参考文档','1. 数组与切片 [https://zhuanlan.zhihu.com/p/456844080](https://zhuanlan.zhihu.com/p/456844080)2. map [https://blog.csdn.net/qq_39312146/article/details/127914811](https://blog.csdn.net/qq_39312146/article/details/127914811)'),(38,'2024-12-17 20:37:34.886','2024-12-17 20:37:34.886',5,'5.判断语句',''),(39,'2024-12-17 20:37:34.886','2024-12-17 20:37:34.886',5,'if语句','以年龄为例，输入的年龄在某一个区间，就输出对应的提示信息```Go<=0     未出生1-18    未成年18-35   青年>=35    中年```很明显，这是一个多选一的情况我们有很多中方式来实现中断式```Gopackage mainimport \"fmt\"func main() {  fmt.Println(\"请输入你的年龄：\")  var age int  fmt.Scan(&age)  if age <= 0 {    fmt.Println(\"未出生\")    return  }  if age <= 18 {    fmt.Println(\"未成年\")    return  }  if age <= 35 {    fmt.Println(\"青年\")    return  }  fmt.Println(\"中年\")}```它也有个好听的名字，叫卫语句嵌套式```Gopackage mainimport \"fmt\"func main() {  fmt.Println(\"请输入你的年龄：\")  var age int  fmt.Scan(&age)  if age <= 18 {    if age <= 0 {      fmt.Println(\"未出生\")    } else {      fmt.Println(\"未成年\")    }  } else {    if age <= 35 {      fmt.Println(\"青年\")    } else {      fmt.Println(\"中年\")    }  }}```多条件式```Gopackage mainimport \"fmt\"func main() {  fmt.Println(\"请输入你的年龄：\")  var age int  fmt.Scan(&age)  if age <= 0 {    fmt.Println(\"未出生\")  }  if age > 0 && age <= 18 {    fmt.Println(\"未成年\")  }  if age > 18 && age <= 35 {    fmt.Println(\"青年\")  }  if age > 35 {    fmt.Println(\"中年\")  }}```'),(40,'2024-12-17 20:37:34.886','2024-12-17 20:37:34.886',5,'switch语句','还是上面那个案例，如果是用switch就很直观了```Gopackage mainimport \"fmt\"func main() {  fmt.Println(\"请输入你的年龄：\")  var age int  fmt.Scan(&age)  switch {  case age <= 0:    fmt.Println(\"未出生\")  case age <= 18:    fmt.Println(\"未成年\")  case age <= 35:    fmt.Println(\"青年\")  default:    fmt.Println(\"中年\")  }}```除了这样的写法，还有枚举所有的可能值例如```Go1   星期一2   星期二3   星期三``````Gopackage mainimport \"fmt\"func main() {  fmt.Println(\"请输入星期数字：\")  var week int  fmt.Scan(&week)  switch week {  case 1:    fmt.Println(\"周一\")  case 2:    fmt.Println(\"周二\")  case 3:    fmt.Println(\"周三\")  case 4:    fmt.Println(\"周四\")  case 5:    fmt.Println(\"周五\")  case 6, 7:    fmt.Println(\"周末\")  default:    fmt.Println(\"错误\")  }}```可以理解为case的值就是switch的枚举结果一般来说，go的switch的多选一，满足其中一个结果之后，就结束switch了例如```Gopackage mainimport \"fmt\"func main() {  fmt.Println(\"请输入你的年龄：\")  var age int  fmt.Scan(&age)  switch {  case age <= 0:    fmt.Println(\"未出生\")  case age <= 18:    fmt.Println(\"未成年\")  case age <= 35:    fmt.Println(\"青年\")  default:    fmt.Println(\"中年\")  }}```我输入一个12，我希望它能输出满足的所有条件，例如我希望它输出，未成年，青年```Gopackage mainimport \"fmt\"func main() {  fmt.Println(\"请输入你的年龄：\")  var age int  fmt.Scan(&age)  switch {  case age <= 0:    fmt.Println(\"未出生\")    fallthrough  case age <= 18:    fmt.Println(\"未成年\")    fallthrough  case age <= 35:    fmt.Println(\"青年\")  default:    fmt.Println(\"中年\")  }}```'),(41,'2024-12-17 20:37:34.886','2024-12-17 20:37:34.886',5,'参考文档','卫语句  [https://blog.csdn.net/wangpaiblog/article/details/114909737](https://blog.csdn.net/wangpaiblog/article/details/114909737)switch [https://blog.csdn.net/qq_43470538/article/details/130436381](https://blog.csdn.net/qq_43470538/article/details/130436381)switch [https://www.php.cn/faq/540744.html](https://www.php.cn/faq/540744.html)');
/*!40000 ALTER TABLE `text_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_article_collect_models`
--

DROP TABLE IF EXISTS `user_article_collect_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_article_collect_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `article_id` bigint(20) unsigned DEFAULT NULL,
  `collect_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`user_id`,`article_id`,`collect_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_article_collect_models`
--

LOCK TABLES `user_article_collect_models` WRITE;
/*!40000 ALTER TABLE `user_article_collect_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_article_collect_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_article_look_history_models`
--

DROP TABLE IF EXISTS `user_article_look_history_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_article_look_history_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `article_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_article_look_history_models`
--

LOCK TABLES `user_article_look_history_models` WRITE;
/*!40000 ALTER TABLE `user_article_look_history_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_article_look_history_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_chat_action_models`
--

DROP TABLE IF EXISTS `user_chat_action_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_chat_action_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `chat_id` bigint(20) unsigned DEFAULT NULL,
  `is_read` tinyint(1) DEFAULT NULL,
  `is_delete` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_chat_action_models`
--

LOCK TABLES `user_chat_action_models` WRITE;
/*!40000 ALTER TABLE `user_chat_action_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_chat_action_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_conf_models`
--

DROP TABLE IF EXISTS `user_conf_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_conf_models` (
  `user_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `like_tags` longtext,
  `update_username_date` datetime(3) DEFAULT NULL,
  `open_collect` tinyint(1) DEFAULT NULL,
  `open_follow` tinyint(1) DEFAULT NULL,
  `open_fans` tinyint(1) DEFAULT NULL,
  `home_style_id` bigint(20) unsigned DEFAULT NULL,
  `look_count` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `uni_user_conf_models_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_conf_models`
--

LOCK TABLES `user_conf_models` WRITE;
/*!40000 ALTER TABLE `user_conf_models` DISABLE KEYS */;
INSERT INTO `user_conf_models` VALUES (1,NULL,NULL,1,1,1,1,0),(2,NULL,NULL,1,1,1,1,0);
/*!40000 ALTER TABLE `user_conf_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_focus_models`
--

DROP TABLE IF EXISTS `user_focus_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_focus_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `focus_user_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_focus_models`
--

LOCK TABLES `user_focus_models` WRITE;
/*!40000 ALTER TABLE `user_focus_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_focus_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_global_notification_models`
--

DROP TABLE IF EXISTS `user_global_notification_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_global_notification_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `notification_id` bigint(20) unsigned DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `is_read` tinyint(1) DEFAULT NULL,
  `is_delete` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_global_notification_models`
--

LOCK TABLES `user_global_notification_models` WRITE;
/*!40000 ALTER TABLE `user_global_notification_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_global_notification_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_login_models`
--

DROP TABLE IF EXISTS `user_login_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_login_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `ip` varchar(32) DEFAULT NULL,
  `addr` varchar(64) DEFAULT NULL,
  `ua` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_login_models`
--

LOCK TABLES `user_login_models` WRITE;
/*!40000 ALTER TABLE `user_login_models` DISABLE KEYS */;
INSERT INTO `user_login_models` VALUES (1,'2024-12-17 20:28:26.405','2024-12-17 20:28:26.405',1,'127.0.0.1','内网','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36'),(2,'2024-12-17 20:29:52.646','2024-12-17 20:29:52.646',1,'127.0.0.1','内网','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36'),(3,'2024-12-17 20:32:44.957','2024-12-17 20:32:44.957',2,'127.0.0.1','内网','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36'),(4,'2024-12-17 20:33:23.568','2024-12-17 20:33:23.568',2,'127.0.0.1','内网','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36'),(5,'2024-12-17 20:38:25.517','2024-12-17 20:38:25.517',1,'127.0.0.1','内网','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36');
/*!40000 ALTER TABLE `user_login_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_message_conf_models`
--

DROP TABLE IF EXISTS `user_message_conf_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_message_conf_models` (
  `user_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `open_comment_message` tinyint(1) DEFAULT NULL,
  `open_digg_message` tinyint(1) DEFAULT NULL,
  `open_private_chat` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `uni_user_message_conf_models_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_message_conf_models`
--

LOCK TABLES `user_message_conf_models` WRITE;
/*!40000 ALTER TABLE `user_message_conf_models` DISABLE KEYS */;
INSERT INTO `user_message_conf_models` VALUES (1,1,1,1),(2,1,1,1);
/*!40000 ALTER TABLE `user_message_conf_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_models`
--

DROP TABLE IF EXISTS `user_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `username` varchar(32) DEFAULT NULL,
  `nickname` varchar(32) DEFAULT NULL,
  `avatar` varchar(256) DEFAULT NULL,
  `abstract` varchar(256) DEFAULT NULL,
  `register_source` tinyint(4) DEFAULT NULL,
  `password` varchar(64) DEFAULT NULL,
  `email` varchar(256) DEFAULT NULL,
  `open_id` varchar(64) DEFAULT NULL,
  `role` tinyint(4) DEFAULT NULL,
  `ip` longtext,
  `addr` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_models`
--

LOCK TABLES `user_models` WRITE;
/*!40000 ALTER TABLE `user_models` DISABLE KEYS */;
INSERT INTO `user_models` VALUES (1,'2024-12-17 20:27:31.549','2024-12-17 20:28:58.167','admin','admin','/uploads/images001/62cc6c0d334c14e239294c14cc36684a.png','',3,'$2a$10$AtNmvRCPsXdtbTjZELvvGu1aqU32Gguvu.52zxFOdsmNp1R9qjaUi','','',1,'',''),(2,'2024-12-17 20:32:34.357','2024-12-17 20:33:14.048','zhangsan','zhangsan','/uploads/images001/06c11c28d0d4488e7d7cb51e9506aedf.png','',3,'$2a$10$BjNrXq/y67KMtycpU/QzJemwlVAVAjEFGGKv3S/wv8J3B/KkejGUW','','',2,'','');
/*!40000 ALTER TABLE `user_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_top_article_models`
--

DROP TABLE IF EXISTS `user_top_article_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_top_article_models` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `article_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`user_id`,`article_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_top_article_models`
--

LOCK TABLES `user_top_article_models` WRITE;
/*!40000 ALTER TABLE `user_top_article_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_top_article_models` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-12-17 21:07:14
