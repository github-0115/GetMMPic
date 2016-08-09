<?php
/* usage:
 * php get_sex_mm.php
 */
 
$url = 'http://1199av.com/html/tupian/siwa/2016/0807/381664.html';
$pic_array = array();
$pic_array_t = array();
$pic_array_pop = array();
$pic_file_name = array();
$pic_array_re_match = array();

$reg_str =  '/<img.*?src=[\"|\']?(.*?)[\"|\']?\s.*?>/i';
$html = file_get_contents($url);

if(preg_match_all($reg_str,$html,$pic_array)) {
	echo "";
} else {
	echo "Not found! ";
}

foreach ($pic_array as $k => $v) {
	if(is_array($v)) {
		foreach($v as $pic_array_t) {
			 if(trim(basename($pic_array_t))  == '>' ) {
				 	//echo "\ntrue\n";
			} else {
				array_push($pic_array_pop, $pic_array_t);
				//array_push($pic_file_name, basename($pic_array_t));
			}
		}
	}
}

function download_image($v) {
	$return_content = file_get_contents($v);
	$filename = basename($v);
	echo "Download file = " . $filename . "\n";
	sleep(1);
	file_put_contents($filename, $return_content);
}

echo "\nGet " . count($pic_array_pop) . "pic \n";
echo "\n Starting download the image\n";
if (is_array($pic_array_pop)) {
	foreach($pic_array_pop as $v) {
		download_image($v);
	}
}
echo "\nFinish download\n";
?>
