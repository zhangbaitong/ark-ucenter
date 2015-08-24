<?php
//error_reporting(0); 

include_once __DIR__ . '/mongodb.php';  
$mongodb=MongoUtil::init();
$db_name=$mongodb->select_db("admin");
//$temp=$mongodb->select_collection($db_name,"db_test");
 $time_start=(int)(microtime(true)*1000);
for($i=0;$i<1000000;$i++)
{
    $obj = array('x' => $i);
    $mongodb->insert("test_table",$obj);
    //$mongodb->save("test_table",$obj);
}
$time_end= (int)(microtime(true)*1000);
echo " use microtime ".($time_end-$time_start)."\r\n";        
//var_dump($temp);

/*
include_once "mongo.php";
$mongo = new HMongodb("127.0.0.1:27017");   
$mongo->selectDb("test_db");   
$mongo->ensureIndex("test_table", array("id"=>1), array('unique'=>true));  
$mongo->insert("test_table", array("id"=>2, "title"=>"asdqw"));  
$mongo->update("test_table", array("id"=>1),array("id"=>1,"title"=>"bbb"));   
$mongo->update("test_table", array("id"=>1),array("id"=>1,"title"=>"bbb"),array("upsert"=>1));   
 */
?>
