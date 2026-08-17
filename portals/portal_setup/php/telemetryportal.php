<?php
session_name('APP_TELE_SESSION');
session_start();

if (!isset($_SESSION['logged_in']) || $_SESSION['logged_in'] !== true) {
    header("Location: telemetry_login.php");
    exit;
}
ini_set('display_errors', 1);
error_reporting(E_ALL);
date_default_timezone_set('UTC');
$directory = __DIR__;
$files = glob($directory . '/*.json');

if (isset($_GET['ajax'])) {
    $entries = isset($_GET['length']) ? (int)$_GET['length'] : 10;
    $offset = isset($_GET['start']) ? (int)$_GET['start'] : 0;

    $searchMac = $_GET['mac'] ?? '';
    $searchDate = $_GET['date'] ?? '';
    $fromTime = $_GET['fromTime'] ?? '';
    $toTime = $_GET['toTime'] ?? '';

    $filtered = [];

    foreach ($files as $file) {
        $filename = basename($file);
        $timestamp = filemtime($file);
        $time = gmdate("Y-m-d H:i:s", $timestamp);

        // Filter MAC address
        if ($searchMac && stripos($filename, $searchMac) === false) continue;

        // Filter date
        if ($searchDate) {
            $fileDate = gmdate('Y-m-d', $timestamp);
            if ($fileDate !== $searchDate) continue;
        }

        // Filter time range
        if ($searchDate && $fromTime && $toTime) {
            $fromTS = strtotime("$searchDate $fromTime UTC");
            $toTS = strtotime("$searchDate $toTime UTC");
            if ($fromTS > $toTS || $timestamp < $fromTS || $timestamp > $toTS) continue;
        }

        $filtered[] = ['filename' => $filename, 'time' => $time];
    }

    // Sort descending by file time
    usort($filtered, fn($a, $b) => strtotime($b['time']) - strtotime($a['time']));

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
    <title>Telemetry Portal</title>
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
        #jsonFiles td:first-child,
	#jsonFiles th:first-child {
  		font-family: 'JetBrains Mono', monospace;
	}
        #jsonFiles_paginate {
            text-align: right !important;
	}
	h2{
             text-align:center;
	}
        #jsonFiles_wrapper {
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
    <div class="rightalign">
     <h5><?php echo htmlspecialchars($_SESSION['username']); ?></h5>
    <a href="telemetry_logout.php">Logout</a>
    </div>
    <h2>Telemetry Portal</h2><hr>
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
      <label for="fromTime">From Time:</label><input type="time" id="fromTime">
      <label for="toTime">To Time:</label><input type="time" id="toTime">
      <button type="button" id="searchBtn">Search</button>
      <button type="button" id="clearBtn">Clear</button>
    </div>
  </div>
</div>


   <table id="jsonFiles" class="display" style="width:100%">
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

        let table = $('#jsonFiles').DataTable({
            processing: true,
            serverSide: true,
	    pageLength: 10,
            lengthChange: false,
            ajax: {
                url: 'telemetryportal.php?ajax=1',
                data: function (d) {
                    d.mac = $('#macSearch').val().trim();
                    d.date = $('#searchDate').val();
                    d.fromTime = $('#fromTime').val();
                    d.toTime = $('#toTime').val();
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
            const fromTime = $('#fromTime').val();
            const toTime = $('#toTime').val();

            if (mac && !isValidMac(mac)) {
                alert("Invalid MAC address format.Please enter a  valid mac address format");
                return;
	    }
	    if (!date && (fromTime || toTime)) {
        	alert("Please select a Date before specifying From or To Time.");
     	   	return;
    	    }

            if (fromTime && toTime && fromTime >= toTime) {
                alert("From Time must be less than To Time.");
                return;
            }

            table.ajax.reload();
        });

        $('#clearBtn').on('click', function () {
            $('#macSearch, #searchDate, #fromTime, #toTime').val('');
            $('#entriesSelect').val('10').trigger('change');
            table.page.len(10).draw();
        });
    </script>
</body>
</html>

