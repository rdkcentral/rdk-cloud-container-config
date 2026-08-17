<?php
 #header('Content-Type: text/plain; charset=utf-8');
 header('Content-Type: application/json;');

/*Global variable */
$tele_path='/var/www/html/logTelemetryServer/upload/';
$debug=0;
//print("Begin Here....");


function getValfromJSON($jsonStr1,$pattern)
{
	$value='';
	foreach($jsonStr1 as $arrIdx =>$jsonStr) {
		if(key($jsonStr)==$pattern) {
			if($debug) print("Matched: ".$jsonStr[$pattern]);
			if($pattern=="mac")
			  $value = str_replace(array(':',' '), '', $jsonStr[$pattern]);
			else
			  $value = $jsonStr[$pattern];
		}
	}
	return $value;
}

try {
	$rtlStr = file_get_contents('php://input');
	$data = json_decode($rtlStr,true);


	if (isset($data['Report']) && is_array($data['Report'])) {
		foreach ($data['Report'] as $item) {
 		   if (isset($item['mac'])) {
       			 $mac_addr = $item['mac'];
    		   }
    	           if (isset($item['Time'])) {
        		$time_raw = $item['Time'];
    		    }
		}

	}
	if (isset($data['searchResult']) && is_array($data['searchResult'])) {
		 foreach ($data['searchResult'] as $item) {
                   if (isset($item['mac'])) {
                         $mac_addr = $item['mac'];
                   }
                   if (isset($item['Time'])) {
                        $time_raw = $item['Time'];
                    }
                }

	}

	$time_str = new DateTime($time_raw);
	$time_conv = $time_str->format('m-d-Y-h-iA');
          $TwoDigitRandomNumber = rand(10,99);
	$json_file=$tele_path.$mac_addr.'_TELE_TEST'.$TwoDigitRandomNumber.$time_conv.'.json';
	if($debug) print($json_file);
if (1) {
	$fp = fopen($json_file, 'w');
	fwrite($fp, $rtlStr."\n");
	fclose($fp);
}
} catch (RuntimeException $e) {
	echo $e->getMessage();
}
?>
