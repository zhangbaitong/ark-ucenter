<?php  
  
final class Config {  
  
    /** 
     * Memcache 配置 
     * @var mixed Memcache 配置 
     */  
    public static $MEMCACHE_SERVER = array(  
        hosts => array(  
            array(  
                host => '127.0.0.1',  
                port => 11211,  
                weight => 100  
            )  
        ),  
        threshold => 30000,  
        saving => 0.2  
    );  
  
    /** 
     * MongoDB 配置 
     * @var mixed  MongoDB 配置  
     */  
    public static $MONGO_SERVER = array(  
        hosts => array(  
            array(ip => '127.0.0.1', port => 27017, username => 'admin', password => 'admin', is_master => true)  
//            ,array(ip => '127.0.0.1', port => 27018, username => 'admin', password => 'admin', is_master => true)  
        )  
    );  
  
    /** 
     * MySQL 配置 
     * @var mixed MySQL 配置 
     */  
    public static $MYSQL_SERVER = array(  
        host=>'127.0.0.1',  
        user=>'root',  
        password=>'123456',  
        port=>3306  
    );  
}  
?>  
