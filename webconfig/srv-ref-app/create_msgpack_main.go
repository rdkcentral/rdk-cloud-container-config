package main

import (
        "bytes"
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "os"
        "regexp"
        "strings"
        "strconv"
        "github.com/vmihailenco/msgpack/v5"
)

type TR181Entry struct {
        Name     string `json:"name" msgpack:"name"`
        Value    string `json:"value" msgpack:"value"`
        DataType int    `json:"dataType" msgpack:"dataType,omitempty"`
}

type TR181Output struct {
        Parameters []TR181Entry `json:"parameters" msgpack:"parameters"`
}

type EmbeddedLanData struct {
  DhcpServerEnable bool `json:"DhcpServerEnable"`
  LanIPAddress string `json:"LanIPAddress"`
  LanSubnetMask string `json:"LanSubnetMask"`
  DhcpStartIPAddress string `json:"DhcpStartIPAddress"`
  DhcpEndIPAddress string `json:"DhcpEndIPAddress"`
  LeaseTime int `json:"LeaseTime"`
}
type EmbeddedLan struct {
  Lan EmbeddedLanData `json:"lan" msgpack:"lan"`
}

type EmbeddedWanData struct {
        Enable bool `json:"Enable"`
        InternalIP string `json:"InternalIP"`
}
type EmbeddedWan struct {
  Wan EmbeddedWanData `json:"wan" msgpack:"wan"`
}

type PortForwardingData struct {
    InternalClient         string `json:"InternalClient"`
    ExternalPortEndRange   string `json:"ExternalPortEndRange"`
    Enable                 bool `json:"Enable"`
    Protocol               string `json:"Protocol"`
    Description            string `json:"Description"`
    ExternalPort           string `json:"ExternalPort"`
}

type PortForwarding struct {
    PortForwarding []PortForwardingData `json:"portforwarding" msgpack:"portforwarding"`
}

type Config struct {
    WebconfigServer struct {
        Host        string `json:"host"`
        Port        int    `json:"port"`
        APIEndpoint string `json:"api_endpoint"`
    } `json:"webconfig_server"`
    FilePermissions string `json:"file_permissions"`
    Logging struct {
        Level  string `json:"level"`
        Output string `json:"output"`
    } `json:"logging"`
    TimeoutSeconds int `json:"timeout_seconds"`
    MaxRetries     int `json:"max_retries"`
}

func loadConfig(configFile string) (*Config, error) {
    data, err := os.ReadFile(configFile)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to parse config file: %w", err)
    }

    return &cfg, nil
}

// Helper to validate IP format
func isValidIP(ip string) bool {
    parts := strings.Split(ip, ".")
    if len(parts) != 4 {
        return false
    }
    for _, part := range parts {
        num, err := strconv.Atoi(part)
        if err != nil || num < 0 || num > 255 {
            return false
        }
    }
    return true
}

// Helper function to parse boolean values
func parseBool(value interface{}) (bool, error) {
    switch v := value.(type) {
    case bool:
        return v, nil
    case string:
        // Handle case-insensitive: "true", "TRUE", "True", "1"
        normalized := strings.ToLower(v)
        switch normalized {
        case "true", "1", "yes":
            return true, nil
        case "false", "0", "no":
            return false, nil
        default:
            return false, fmt.Errorf("invalid boolean value: %q", v)
        }
    case float64:
        if v == 0 {
            return false, nil
        } else if v == 1 {
            return true, nil
        }
        return false, fmt.Errorf("invalid boolean numeric value: %v", v)
    case nil:
        return false, fmt.Errorf("boolean value is nil")
    default:
        return false, fmt.Errorf("unexpected type for boolean value: %T, value: %v", v, v)
    }
}

// Helper function to parse integer values
func parseInt(value interface{}) (int, error) {
    switch v := value.(type) {
    case float64:
        return int(v), nil
    case int:
        return v, nil
    case string:
        // Try to parse string as integer
        intVal, err := strconv.Atoi(v)
        if err != nil {
            return 0, fmt.Errorf("cannot convert string to int: %q - %v", v, err)
        }
        return intVal, nil
    case nil:
        return 0, fmt.Errorf("integer value is nil")
    default:
        return 0, fmt.Errorf("unexpected type for integer value: %T, value: %v", v, v)
    }
}

// SafeGetString with default value
func safeGetString(data map[string]interface{}, key string, defaultVal string) (string, error) {
    val, exists := data[key]
    if !exists {
        return defaultVal, fmt.Errorf("key %q not found in data", key)
    }

    if val == nil {
        return defaultVal, fmt.Errorf("key %q has nil value", key)
    }

    strVal, ok := val.(string)
    if !ok {
        return defaultVal, fmt.Errorf("key %q expected string, got %T", key, val)
    }

    if strVal == "" {
        return defaultVal, fmt.Errorf("key %q is empty string", key)
    }

    return strVal, nil
}

// SafeGetInt with default value
func safeGetInt(data map[string]interface{}, key string, defaultVal int) (int, error) {
    val, exists := data[key]
    if !exists {
        return defaultVal, fmt.Errorf("key %q not found in data", key)
    }

    if val == nil {
        return defaultVal, fmt.Errorf("key %q has nil value", key)
    }

    return parseInt(val)
}

// SafeGetBool with default value
func safeGetBool(data map[string]interface{}, key string, defaultVal bool) (bool, error) {
    val, exists := data[key]
    if !exists {
        return defaultVal, fmt.Errorf("key %q not found in data", key)
    }

    if val == nil {
        return defaultVal, fmt.Errorf("key %q has nil value", key)
    }

    return parseBool(val)
}

func handleLan(subdocData map[string]interface{}) (EmbeddedLan, error) {
    dhcp, err := safeGetBool(subdocData, "DhcpServerEnable", false)
    if err != nil {
        return EmbeddedLan{}, err
    }

    lanIP, err := safeGetString(subdocData, "LanIPAddress", "")
    if err != nil {
        return EmbeddedLan{}, err
    }

    lanSubnet, err := safeGetString(subdocData, "LanSubnetMask", "")
    if err != nil {
        return EmbeddedLan{}, err
    }

    dhcpStart, err := safeGetString(subdocData, "DhcpStartIPAddress", "")
    if err != nil {
        return EmbeddedLan{}, err
    }

    dhcpEnd, err := safeGetString(subdocData, "DhcpEndIPAddress", "")
    if err != nil {
        return EmbeddedLan{}, err
    }

    leaseTime, err := safeGetInt(subdocData, "LeaseTime", 3600)
    if err != nil {
        return EmbeddedLan{}, err
    }


    return EmbeddedLan{
        Lan: EmbeddedLanData{
            DhcpServerEnable:   dhcp,
            LanIPAddress:       lanIP,
            LanSubnetMask:      lanSubnet,
            DhcpStartIPAddress: dhcpStart,
            DhcpEndIPAddress:   dhcpEnd,
            LeaseTime:          leaseTime,
        },
    }, nil
}

func handleWan(subdocData map[string]interface{}) (EmbeddedWan, error) {
    enable, err := safeGetBool(subdocData, "Enable", false)
    if err != nil {
        return EmbeddedWan{}, err
    }

    internalIP, err := safeGetString(subdocData, "InternalIP", "")
    if err != nil {
        return EmbeddedWan{}, err
    }

    if !isValidIP(internalIP) {
        return EmbeddedWan{}, fmt.Errorf("invalid InternalIP: %q", internalIP)
    }

    return EmbeddedWan{
        Wan: EmbeddedWanData{
            Enable:     enable,
            InternalIP: internalIP,
        },
    }, nil
}

func handlePortForwarding(subdocData interface{}) (PortForwarding, error) {
    subdocList, ok := subdocData.([]interface{})
    if !ok {
        return PortForwarding{}, fmt.Errorf("portforwarding data is not a list, got %T", subdocData)
    }

    var portForwardingList []PortForwardingData

    for idx, item := range subdocList {
        itemMap, ok := item.(map[string]interface{})
        if !ok {
            return PortForwarding{}, fmt.Errorf("portforwarding item %d is not a map, got %T", idx, item)
        }

        internalClient, err := safeGetString(itemMap, "InternalClient", "")
        if err != nil {
            return PortForwarding{}, fmt.Errorf("item %d: %w", idx, err)
        }

        externalPortEndRange, err := safeGetString(itemMap, "ExternalPortEndRange", "")
        if err != nil {
            return PortForwarding{}, fmt.Errorf("item %d: %w", idx, err)
        }

        enable, err := safeGetBool(itemMap, "Enable", false)
        if err != nil {
            return PortForwarding{}, fmt.Errorf("item %d: %w", idx, err)
        }

        protocol, err := safeGetString(itemMap, "Protocol", "")
        if err != nil {
            return PortForwarding{}, fmt.Errorf("item %d: %w", idx, err)
        }

        description, err := safeGetString(itemMap, "Description", "")
        if err != nil {
            return PortForwarding{}, fmt.Errorf("item %d: %w", idx, err)
        }

        externalPort, err := safeGetString(itemMap, "ExternalPort", "")
        if err != nil {
            return PortForwarding{}, fmt.Errorf("item %d: %w", idx, err)
        }

        portForwardingList = append(portForwardingList, PortForwardingData{
            InternalClient:       internalClient,
            ExternalPortEndRange: externalPortEndRange,
            Enable:               enable,
            Protocol:             protocol,
            Description:          description,
            ExternalPort:         externalPort,
        })
    }

    return PortForwarding{
        PortForwarding: portForwardingList,
    }, nil
}

func unmarshalJSON(data []byte) (interface{}, error) {
    var result interface{}
    if err := json.Unmarshal(data, &result); err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %v", err)
    }
    return result, nil
}

func createSubdoc(jsonFile string, subdocType string, tr181 string) error {
    var subdoc interface{}
    var err error

    // Read JSON file
    data, err := os.ReadFile(jsonFile)
    if err != nil {
        return fmt.Errorf("failed to read JSON file: %w", err)
    }

    // Unmarshal JSON into map[string]interface{}
    var subdocData map[string]interface{}
    if subdocType == "portforwarding" {
        // Parse JSON data as a slice for portforwarding, as its a list
        var subdocDataList []interface{}
        if err := json.Unmarshal(data, &subdocDataList); err != nil {
            log.Fatalf("Failed to parse JSON as a list: %v", err)
            return err
        }
        subdocData = map[string]interface{}{"listData": subdocDataList}
    } else {
        if err := json.Unmarshal(data, &subdocData); err != nil {
            log.Fatalf("Failed to parse JSON: %v", err)
            return err
        }
    }

   // Handle subdoc types
    switch subdocType {
    case "lan":
        subdoc, err = handleLan(subdocData)
    case "wan":
        subdoc, err = handleWan(subdocData)
    case "portforwarding":
        subdoc, err = handlePortForwarding(subdocData["listData"])
    case "privatessid":
        subdoc = subdocData
    default:
        subdoc = subdocData
    }

    if err != nil {
        return fmt.Errorf("failed to build subdoc for %q: %w", subdocType, err)
    }

    prettyJSON, err := json.MarshalIndent(subdoc, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal JSON: %w", err)
    }
    if os.Getenv("DEBUG") == "1" {
        fmt.Println(string(prettyJSON))
    }

    // Serialize subdoc into MsgPack format
    stringifiedBin, err := msgpack.Marshal(subdoc)
    if err != nil {
        log.Fatalf("Failed to encode msgpack: %v", err)
        return err
    }

    // Create final blob
    blobData := TR181Output{
            Parameters: []TR181Entry{
                    {
                            Name:     tr181,
                            Value:    string(stringifiedBin),
                            DataType: 12,
                    },
            },
    }

    // MsgPack final encoding
    msgpackFormat, err := msgpack.Marshal(blobData)
    if err != nil {
            log.Fatalf("Failed to encode blob msgpack: %v", err)
            return err
    }

    // Write to file
    msgpackFile := fmt.Sprintf("%s.msgpack", subdocType)
    if err := os.WriteFile(msgpackFile, msgpackFormat, 0644); err != nil {
            log.Fatalf("Failed to create blob msgpack file: %v", err)
            return err
    }

    return nil
}

func validateMACAddress(macAddr string) error {
    // Accept only uppercase hex without separators, e.g. A1B2C3D4E5F6.
   // if !regexp.MustCompile(`^[0-9A-F]{12}$`).MatchString(macAddr) {
    macRegex := regexp.MustCompile(`^(?:[0-9A-Fa-f]{2}[:-]){5}(?:[0-9A-Fa-f]{2})|(?:[0-9A-Fa-f]{2}){6}$`)

    if !macRegex.MatchString(macAddr) {
        return fmt.Errorf("invalid MAC address format: %q", macAddr)
    }
    return nil
}

func writeToDB(macAddr string, subdocName string, cfg *Config) error {
    // Validate MAC address first
    if err := validateMACAddress(macAddr); err != nil {
        return err
    }

    // Build URL from config
    url := fmt.Sprintf(
        "http://%s:%d%s",
        cfg.WebconfigServer.Host,
        cfg.WebconfigServer.Port,
        fmt.Sprintf(cfg.WebconfigServer.APIEndpoint, macAddr, subdocName),
    )

    // Read MsgPack file
    msgpackBin := fmt.Sprintf("%s.msgpack", subdocName)
    fileData, err := os.ReadFile(msgpackBin)
    if err != nil {
        return fmt.Errorf("failed to read msgpack file: %w", err)
    }

    // Send POST request
    req, err := http.NewRequest(
        "POST",
	url,
	bytes.NewReader(fileData),
//	io.NopCloser(bytes.NewReader(fileData)),
    )
    if err != nil {
        return fmt.Errorf("failed to create POST request: %v", err)
    }
    req.Header.Set("Content-Type", "application/msgpack")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send POST request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("request failed with HTTP response: %d", resp.StatusCode)
    }

    fmt.Println("Request successful")
    return nil
}

func main() {
    if len(os.Args) < 5 {
        fmt.Println("Usage: go run create_msgpack_main.go <subdoc.json> <subdoc_name> <tr181> <MAC>")
        os.Exit(1)
    }

    subdocJSON := os.Args[1]
    subdocName := os.Args[2]
    tr181 := os.Args[3]
    macAddr := os.Args[4]

    // Load config
    cfg, err := loadConfig("config.json")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
        os.Exit(1)
    }

    // Create MsgPack binary from JSON data
    if err := createSubdoc(subdocJSON, subdocName, tr181); err != nil {
        fmt.Fprintf(os.Stderr, "Error creating MsgPack binary: %v\n", err)
        os.Exit(1)
    }

    // Push MsgPack data to the server
    if err := writeToDB(macAddr, subdocName, cfg); err != nil {
        fmt.Fprintf(os.Stderr, "Error pushing MsgPack data: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("MsgPack creation and upload successful")
}

