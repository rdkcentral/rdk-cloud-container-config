<?php
 header('Content-Type: text/plain; charset=utf-8');

 /*Global variable */
 $fsize_lim=(10*1024*1024); //10MB

 try {
 	if ( !isset($_FILES['filename']['error']) || is_array($_FILES['filename']['error']) ) {
		print ($_FILES['filename']['error']);
		throw new RuntimeException('Invalid parameters.');
	}

	// Check $_FILES['filename']['error'] value.
	switch ($_FILES['filename']['error']) {
		case UPLOAD_ERR_OK:
			break;
		case UPLOAD_ERR_NO_FILE:
			throw new RuntimeException('No file sent.');
		case UPLOAD_ERR_INI_SIZE:
		case UPLOAD_ERR_FORM_SIZE:
			throw new RuntimeException('Exceeded permitted filesize.');
		default:
			throw new RuntimeException('Unknown error.');
	}
	$target_file=$_FILES['filename']['name'];
	$tmp_filename=$_FILES['filename']['tmp_name'];

	if ($_FILES['filename']['size'] > $fsize_lim) {
		throw new RuntimeException('Exceeded filesize limit.');
	}

	// Allow tar and text files only.
	$finfo = new finfo(FILEINFO_MIME_TYPE);
	print("info ".$finfo->file($tmp_filename));
	if (false === $ext = array_search( $finfo->file($tmp_filename),
				array(
				'tgz' => 'application/gzip',
				'txt' => 'text/plain',
				),
				true
	)) {
		throw new RuntimeException('Invalid file format: '.$ext);
	}

	$file_suffix = strtolower(pathinfo($target_file,PATHINFO_EXTENSION));
	if($file_suffix != "log" && $file_suffix != "txt" && $file_suffix != "tgz" && $file_suffix != "zip") {
		throw new RuntimeException('Invalid file extension: [.'.$file_suffix.'] (Allowed .tgz,.txt)');
	}

	if ( !move_uploaded_file( $tmp_filename, 
		sprintf('./upload/%s', $target_file) ) ) {
		throw new RuntimeException('Failed to move uploaded file.');
	}

	echo 'File is uploaded successfully.';

 } catch (RuntimeException $e) {
	echo $e->getMessage();
 }

?>
