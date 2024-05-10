## 37.38

### 单表查询
前期表准备
create table emp(
  id int not null unique auto_increment,
  name varchar(20) not null,
  sex enum('male','female') not null default 'male', #大部分是男的
  age int(3) unsigned not null default 28,
  hire_date date not null,
  post varchar(50),
  post_comment varchar(100),
  salary double(15,2),
  office int, #一个部门一个屋子
  depart_id int
) charset gbk;

#插入记录
#三个部门：教学，销售，运营
insert into emp(name,sex,age,hire_date,post,salary,office,depart_id) values
('jason','male',18,'20170301',"张江第一帅形象代言",7300.33,401,1), #以下是教学部
('egon','male',78,'20150302','teacher',1000000.31,401,1),
('kevin','male',81,'20130305','teacher',8300,401,1),
('tank','male',73,'20140701','teacher',3500,401,1),
('owen','male',28,'20121101','teacher',2100,401,1),
('jerry','female',18,'20110211','teacher',9000,401,1),
('nick','male',18,'19000301','teacher',30000,401,1),
('sean','male',48,'20101111','teacher',10000,401,1),

('歪歪','female',48,'20150311','sale',3000.13,402,2),#以下是销售部门
('丫丫','female',38,'20101101','sale',2000.35,402,2),
('丁丁','female',18,'20110312','sale',1000.37,402,2),
('星星','female',18,'20160513','sale',3000.29,402,2),
('格格','female',28,'20170127','sale',4000.33,402,2),

('张野','male',28,'20160311','operation',10000.13,403,3), #以下是运营部门
('程咬金','male',18,'19970312','operation',20000,403,3),
('程咬银','female',18,'20130311','operation',19000,403,3),
('程咬铜','male',18,'20150411','operation',18000,403,3),
('程咬铁','female',18,'20140512','operation',17000,403,3)
;

#ps：如果在windows系统中，插入中文字符，select的结果为空白，可以将所有字符编码统一设置成gbk
1.语法执行顺序
# 初识查询语句
select id,name from emp where id >= 3 and id <= 6;
# 先后顺序
from
where
select
2.where约束条件
# 1.查询id大于等于3小于等于6的数据
select id,name from emp where id >= 3 and id <= 6;
select *  from emp where id between 3 and 6;  

# 2.查询薪资是20000或者18000或者17000的数据
select * from emp where salary = 20000 or salary = 18000 or salary = 17000;
select * from emp where salary in (20000,18000,17000);  # 简写

# 3.查询员工姓名中包含o字母的员工姓名和薪资
# 在你刚开始接触mysql查询的时候，建议你按照查询的优先级顺序拼写出你的sql语句
"""
先是查哪张表 from emp
再是根据什么条件去查 where name like ‘%o%’
再是对查询出来的数据筛选展示部分 select name,salary
"""
select name,salary from emp where name like '%o%';

# 4.查询员工姓名是由四个字符组成的员工姓名与其薪资
select name,salary from emp where name like '____';
select name,salary from emp where char_length(name) = 4;

# 5.查询id小于3或者大于6的数据
select *  from emp where id not between 3 and 6;

# 6.查询薪资不在20000，18000，17000范围的数据
select * from emp where salary not in (20000,18000,17000);

# 7.查询岗位描述为空的员工名与岗位名  针对null不能用等号，只能用is
select name,post from emp where post_comment = NULL;  # 查询为空！
select name,post from emp where post_comment is NULL;
select name,post from emp where post_comment is not NULL;
3.group by
# 数据分组应用场景：每个部门的平均薪资，男女比例等

# 1.按部门分组
select * from emp group by post;  # 分组后取出的是每个组的第一条数据
select id,name,sex from emp group by post;  # 验证
"""
设置sql_mode为only_full_group_by，意味着以后但凡分组，只能取到分组的依据，
不应该在去取组里面的单个元素的值，那样的话分组就没有意义了，因为不分组就是对单个元素信息的随意获取
"""
set global sql_mode="strict_trans_tables,only_full_group_by";
# 重新链接客户端
select * from emp group by post;  # 报错
select id,name,sex from emp group by post;  # 报错
select post from emp group by post;  # 获取部门信息
# 强调:只要分组了，就不能够再“直接”查找到单个数据信息了，只能获取到组名


# 2.获取每个部门的最高工资  
# 以组为单位统计组内数据>>>聚合查询(聚集到一起合成为一个结果)
# 每个部门的最高工资
select post,max(salary) from emp group by post;
# 每个部门的最低工资
select post,min(salary) from emp group by post;
# 每个部门的平均工资
select post,avg(salary) from emp group by post;
# 每个部门的工资总和
select post,sum(salary) from emp group by post;
# 每个部门的人数
select post,count(id) from emp group by post;

# 3.查询分组之后的部门名称和每个部门下所有的学生姓名
# group_concat（分组之后用）不仅可以用来显示除分组外字段还有拼接字符串的作用
select post,group_concat(name) from emp group by post;

select post,group_concat(name,"_SB") from emp group by post;

select post,group_concat(name,": ",salary) from emp group by post;

select post,group_concat(salary) from emp group by post;


# 4.补充concat（不分组时用）拼接字符串达到更好的显示效果 as语法使用
select name as 姓名,salary as 薪资 from emp;
select concat("NAME: ",name) as 姓名,concat("SAL: ",salary) as 薪资 from emp;

# 补充as语法 即可以给字段起别名也可以给表起
select emp.id,emp.name from emp as t1; # 报错  因为表名已经被你改成了t1
select t1.id,t1.name from emp as t1;

# 查询四则运算
# 查询每个人的年薪
select name,salary*12 as annual_salary from emp;
select name,salary*12 annual_salary from emp;  # as可以省略
练习题
# 刚开始查询表，一定要按照最基本的步骤，先确定是哪张表，再确定查这张表也没有限制条件，再确定是否需要分类，最后再确定需要什么字段对应的信息

1. 查询岗位名以及岗位包含的所有员工名字
2. 查询岗位名以及各岗位内包含的员工个数
3. 查询公司内男员工和女员工的个数
4. 查询岗位名以及各岗位的平均薪资
5. 查询岗位名以及各岗位的最高薪资
6. 查询岗位名以及各岗位的最低薪资
7. 查询男员工与男员工的平均薪资，女员工与女员工的平均薪资
"""
参考答案：
select post,group_concat(name) from emp group by post;
select post,count(id) from emp group by post;
select sex,count(id) from employee group by sex;
select post,avg(salary) from emp group by post;
select post,max(salary) from employee group by post;
select post,min(salary) from employee group by post;
select sex,avg(salary) from employee group by sex;
"""

# 关键字where group by同时出现的情况下，group by必须在where之后
# where先对整张表进行一次筛选，如何group by再对筛选过后的表进行分组
# 如何验证where是在group by之前执行而不是之后 利用聚合函数 因为聚合函数只能在分组之后才能使用
select id,name,age from emp where max(salary) > 3000;  # 报错！

select max(salary) from emp;  
# 正常运行，不分组意味着每一个人都是一组，等运行到max(salary)的时候已经经过where，group by操作了，只不过我们都没有写这些条件

# 语法顺序
select
from
where
group by

# 再识执行顺序
from
where 
group by
select


8、统计各部门年龄在30岁以上的员工平均工资
select post,avg(salary) from emp where age > 30 group by post; 
# 对where过滤出来的虚拟表进行一个分组

# 还不明白可以分步执行查看结构
select * from emp where age>30;
# 基于上面的虚拟表进行分组
select * from emp where age>=30 group by post;
4.having
截止目前已经学习的语法
select 查询字段1,查询字段2,... from 表名
        where 过滤条件
        group by分组依据

# 语法这么写，但是执行顺序却不一样
from
where
group by
select
    having的语法格式与where一致，只不过having是在分组之后进行的过滤，即where虽然不能用聚合函数，但是having可以！
1、统计各部门年龄在30岁以上的员工平均工资，并且保留平均工资大于10000的部门
select post,avg(salary) from emp
        where age >= 30
        group by post
        having avg(salary) > 10000;
# 如果不信你可以将having取掉，查看结果，对比即可验证having用法！

#强调：having必须在group by后面使用
select * from emp having avg(salary) > 10000;  # 报错
5.distinct
# 对有重复的展示数据进行去重操作
select distinct post from emp;
6.order by
select * from emp order by salary asc; #默认升序排
select * from emp order by salary desc; #降序排

select * from emp order by age desc; #降序排

#先按照age降序排，在年轻相同的情况下再按照薪资升序排
select * from emp order by age desc,salary asc; 

# 统计各部门年龄在10岁以上的员工平均工资，并且保留平均工资大于1000的部门，然后对平均工资进行排序
select post,avg(salary) from emp
    where age > 10
    group by post
    having avg(salary) > 1000
    order by avg(salary)
    ;
7.limit
# 限制展示条数
select * from emp limit 3;
# 查询工资最高的人的详细信息
select * from emp order by salary desc limit 1;

# 分页显示
select * from emp limit 0,5;  # 第一个参数表示起始位置，第二个参数表示的是条数，不是索引位置
select * from emp limit 5,5;
8.正则
    select * from emp where name regexp '^j.*(n|y)$';
多表查询
表创建
#建表
create table dep({1}
id int,
name varchar(20) 
);

create table emp1({1}
id int primary key auto_increment,
name varchar(20),
sex enum('male','female') not null default 'male',
age int,
dep_id int
);

#插入数据
insert into dep values
(200,'技术'),
(201,'人力资源'),
(202,'销售'),
(203,'运营');

insert into emp(name,sex,age,dep_id) values
('jason','male',18,200),
('egon','female',48,201),
('kevin','male',38,201),
('nick','female',28,202),
('owen','male',18,200),
('jerry','female',18,204)
{1};

# 当初为什么我们要分表，就是为了方便管理，在硬盘上确实是多张表，但是到了内存中我们应该把他们再拼成一张表进行查询才合理
表查询
select * from emp,dep;  # 左表一条记录与右表所有记录都对应一遍>>>笛卡尔积

# 将所有的数据都对应了一遍，虽然不合理但是其中有合理的数据，现在我们需要做的就是找出合理的数据

# 查询员工及所在部门的信息
select * from emp,dep where emp.dep_id = dep.id;
# 查询部门为技术部的员工及部门信息
select * from emp,dep where emp.dep_id = dep.id and dep.name = '技术';


# 将两张表关联到一起的操作，有专门对应的方法
# 1、内连接：只取两张表有对应关系的记录
select * from emp inner join dep on emp.dep_id = dep.id;
select * from emp inner join dep on emp.dep_id = dep.id
                            where dep.name = "技术";

# 2、左连接: 在内连接的基础上保留左表没有对应关系的记录
select * from emp left join dep on emp.dep_id = dep.id;

# 3、右连接: 在内连接的基础上保留右表没有对应关系的记录
select * from emp right join dep on emp.dep_id = dep.id;

# 4、全连接：在内连接的基础上保留左、右面表没有对应关系的的记录
select * from emp left join dep on emp.dep_id = dep.id
union
select * from emp right join dep on emp.dep_id = dep.id;
子查询
# 就是将一个查询语句的结果用括号括起来当作另外一个查询语句的条件去用
# 1.查询部门是技术或者人力资源的员工信息
"""
先获取技术部和人力资源部的id号，再去员工表里面根据前面的id筛选出符合要求的员工信息
"""
select * from emp where dep_id in (select id from dep where name = "技术" or name = "人力资源");

# 2.每个部门最新入职的员工 思路：先查每个部门最新入职的员工，再按部门对应上联表查询
select t1.id,t1.name,t1.hire_date,t1.post,t2.* from emp as t1
inner join
(select post,max(hire_date) as max_date from emp group by post) as t2
on t1.post = t2.post
where t1.hire_date = t2.max_date
;

# 3.查询平均年轻在25岁以上的部门名
select name from dep 
			where id in 
			(select dep_id from emp group by dep_id having avg(age)>25);
 
select dep.name from emp inner join dep on emp.dep_id = dep.id 
			group by dep.name
			having avg(age) > 25;
 
# exist(了解)
EXISTS关字键字表示存在。在使用EXISTS关键字时，内层查询语句不返回查询的记录，
而是返回一个真假值，True或False。
当返回True时，外层查询语句将进行查询
当返回值为False时，外层查询语句不进行查询。
select * from employee
    where exists
    (select id from department where id > 3);
 
select * from employee
    where exists
    (select id from department where id > 250);
 
 
"""
记住一个规律，表的查询结果可以作为其他表的查询条件，也可以通过其别名的方式把它作为一张虚拟表去跟其他表做关联查询
"""

select * from emp inner join dep on emp.dep_id = dep.id;