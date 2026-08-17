<?php
session_name("APP_LOG_SESSION");

session_start();

if (!isset($_SESSION['logged_in']) || $_SESSION['logged_in'] !== true) {
    header("Location: logportal_login.php");
    exit;
}

ini_set('display_errors', 1);
error_reporting(E_ALL);
date_default_timezone_set('UTC');
$directory = __DIR__;
$files = glob($directory . '/*.tgz');

if (isset($_GET['ajax'])) {
    $entries = isset($_GET['length']) ? (int)$_GET['length'] : 10;
    $offset = isset($_GET['start']) ? (int)$_GET['start'] : 0;

    $searchMac = $_GET['mac'] ?? '';
    $searchDate = $_GET['date'] ?? '';

    
    $filtered = [];

   foreach ($files as $file) {

    $filename  = basename($file);
    $timestamp = filemtime($file);

    // Filter MAC address
    if ($searchMac && stripos($filename, $searchMac) === false) {
        continue;
    }

    // Filter date (UTC)
    if ($searchDate) {
        $fileDate = gmdate('Y-m-d', $timestamp);
        if ($fileDate !== $searchDate) {
            continue;
        }
    }

    // ✅ MUST be inside the foreach
    $filtered[] = [
        'filename'  => $filename,
        'timestamp' => $timestamp,                     // for sorting
        'time'      => gmdate("Y-m-d H:i:s", $timestamp) // display only
    ];
  }	

	/* -------- SORT (newest first) -------- */
	if (!empty($filtered)) {
    		usort($filtered, fn($a, $b) => $b['timestamp'] <=> $a['timestamp']);
  }	


    $pageData = array_slice($filtered, $offset, $entries);
    $response = [
        'draw' => intval($_GET['draw'] ?? 1),
        'recordsTotal' => count($files),
        'recordsFiltered' => count($filtered),
        'data' => array_map(fn($f) => [
            "<a href='{$f['filename']}'>{$f['filename']}</a>",
            $f['time']
        ], $pageData)
    ];


    header('Content-Type: application/json');
    echo json_encode($response);
    exit;
}
?>

<!DOCTYPE html>
<html>
<head>
    <title>Log Portal</title>
    <link rel="stylesheet" href="https://cdn.datatables.net/1.13.6/css/jquery.dataTables.min.css">
    <style>
        .filters {
            display: flex;
            flex-wrap: wrap;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 10px;
        }
        .filters > div {
            margin: 5px 0;
        }
        label {
            margin-right: 6px;
            font-weight: bold;
        }
        input, select, button {
            margin-right: 10px;
            padding: 5px;
	}
        #tgzfiles td:first-child,
	#tgzfiles th:first-child {
  		font-family: 'JetBrains Mono', monospace;
	}
        #tgzfiles_paginate {
            text-align: right !important;
	}
	h2{
             text-align:center;
	}
        #tgzfiles_wrapper {
 	 max-width: 70%;   /* Reduce to 70% of parent/container */
  	 margin: 0 auto;   /* Center it horizontally */
	}

.filter-bar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 20px; /* space between entries and filters */
  margin-bottom: 15px;
}

.entries-section {
  display: flex;
  align-items: center;
  gap: 5px;
}

.filters-section {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.filters-section label {
  margin-right: 4px;
}

.filters-container {
 flex: 1 1 auto;
  display: flex;
  justify-content: center;
}

.filters-section input[type="text"],
.filters-section input[type="date"],
.filters-section input[type="time"] {
  width: 140px;
}

.rightalign{
   text-align: right;
}


    </style>
</head>
<body>
    <h2>Log Portal</h2>
 <div class="rightalign">
     <h5><?php echo htmlspecialchars($_SESSION['username']); ?></h5>
    <a href="logportal_logout.php">Logout</a>
  </div>
 <hr>
 <div class="filter-bar">
   <div class="entries-section">
    <label for="entriesSelect">Entries:</label>
    <select id="entriesSelect">
      <option value="10" selected>10</option>
      <option value="20">20</option>
      <option value="50">50</option>
      <option value="100">100</option>
    </select>
  </div>

  <div class="filters-container">
    <div class="filters-section">
      <label for="macSearch">MAC:</label><input type="text" id="macSearch" placeholder="AA:BB:CC:DD:EE:FF">
      <label for="searchDate">Date: </label><input type="date" id="searchDate">
      <button type="button" id="searchBtn">Search</button>
      <button type="button" id="clearBtn">Clear</button>
    </div>
  </div>
</div>


   <table id="tgzfiles" class="display" style="width:100%">
        <thead>
            <tr>
                <th>Filename</th>
                <th>Date</th>
            </tr>
        </thead>
    </table> 

    <script src="https://code.jquery.com/jquery-3.7.1.min.js"></script>
    <script src="https://cdn.datatables.net/1.13.6/js/jquery.dataTables.min.js"></script>
    <script>
        function isValidMac(mac) {
            return /^([0-9A-F]{2}:){5}[0-9A-F]{2}$/i.test(mac);
        }

        let table = $('#tgzfiles').DataTable({
            processing: true,
            serverSide: true,
	    pageLength: 10,
            lengthChange: false,
            ajax: {
                url: 'logportal.php?ajax=1',
                data: function (d) {
                    d.mac = $('#macSearch').val().trim();
                    d.date = $('#searchDate').val();
                }
            },
            pagingType: "full_numbers",
	    // language: { search: "" },
	    dom:'lrtip',
        });

        $('#entriesSelect').on('change', function () {
            table.page.len(this.value).draw();
	});

        const today = new Date().toISOString().split('T')[0];
        $('input[type="date"]').attr('max', today);


        $('#searchBtn').on('click', function () {
	    const mac = $('#macSearch').val().trim();
	    const date = document.getElementById('searchDate').value;

            if (mac && !isValidMac(mac)) {
                alert("Invalid MAC address format.Please enter a  valid mac address format");
                return;
	    }

            table.ajax.reload();
        });

        $('#clearBtn').on('click', function () {
            $('#macSearch, #searchDate').val('');
            $('#entriesSelect').val('10').trigger('change');
            table.page.len(10).draw();
        });
    </script>
</body>
</html>

