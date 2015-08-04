
##用户表其他属性（如qq、手机号码、指纹存储在mongdb）
create table IF NOT EXISTS account_tab(
  ac_id bigint(20) not null AUTO_INCREMENT COMMENT '账号ID',
  ac_name varchar(100) not null COMMENT '账户名称',
  ac_password varchar(50) not null COMMENT '账户密码',
  status   integer not null COMMENT '账户状态0:启用;1:停用;',
  source   integer not null DEFAULT 0 COMMENT '账户来源0:直接注册;1:QQ;',
  mid       varchar(32) not null COMMENT 'mongodb object id',
  create_time   integer not null,
  PRIMARY KEY (`ac_id`)
)ENGINE=MyISAM  DEFAULT CHARSET=utf8 COMMENT='认证中心-账号表'  AUTO_INCREMENT=1;

create table IF NOT EXISTS app_info_tab(
  app_id varchar(50) not null  COMMENT '应用ID',
  app_key varchar(100) not null COMMENT '应用密钥',
  app_name varchar(100) not null COMMENT '应用名称',
  app_desc varchar(256) not null COMMENT '应用描述',
  domain  varchar(256) not null COMMENT '域名',
  status   integer(1) not null COMMENT '账户状态0:初始化;1:审核通过;2:审核未通过;3:停用',
  create_time   integer not null,
  PRIMARY KEY (`app_id`)
)ENGINE=MyISAM  DEFAULT CHARSET=utf8 COMMENT='认证中心-应用表'  AUTO_INCREMENT=1  ;

create table IF NOT EXISTS resource_tab(
  res_id bigint(20) not null AUTO_INCREMENT COMMENT '权限ID',
  app_id varchar(50) not null  COMMENT '应用ID',
  res_name  varchar(50) not null COMMENT '资源名称',
  res_cname  varchar(50) not null COMMENT '资源中文名称',
  res_type int(1) not null COMMENT '调用类型 0:http get;1:http post;2:https get;3:https post',
  res_target varchar(512) not null COMMENT '资源目标',
  res_desc varchar(256) not null COMMENT '资源描述',
  status  integer not null COMMENT '状态0:启用;1:停用;',
  create_time  integer not null,
  PRIMARY KEY (`res_id`)
)ENGINE=MyISAM  DEFAULT CHARSET=utf8 COMMENT='资源表'  AUTO_INCREMENT=1;


create table IF NOT EXISTS app_confered_tab(
  app_id varchar(50) not null COMMENT '应用ID',
  res_id bigint(20) not null COMMENT '权限ID',
  status   integer not null COMMENT '状态0:启用;1:停用;',
  create_time   integer not null,
  PRIMARY KEY (`app_id`,`res_id`)
)ENGINE=MyISAM  DEFAULT CHARSET=utf8 COMMENT='应用授权结果表';

create table IF NOT EXISTS app_confered_person_tab(
  app_id varchar(50) not null COMMENT '应用ID',
  openid varchar(50) not null COMMENT '账号ID',
  res_id bigint(20) not null COMMENT '权限ID',
  status   integer not null COMMENT '状态0:启用;1:停用;',
  create_time   integer not null,
  PRIMARY KEY (`app_id`,`openid`,`res_id`)
)ENGINE=MyISAM  DEFAULT CHARSET=utf8 COMMENT='个人应用授权结果表' ;

create table IF NOT EXISTS openid_tab(
  ac_id bigint(20) not null COMMENT '账号ID',
  app_id varchar(50) not null COMMENT '应用ID',
  openid   varchar(50) not null COMMENT 'open id',
  PRIMARY KEY (`app_id`,`ac_id`)
)ENGINE=MyISAM  DEFAULT CHARSET=utf8 COMMENT='账号open_id对应表' ;

