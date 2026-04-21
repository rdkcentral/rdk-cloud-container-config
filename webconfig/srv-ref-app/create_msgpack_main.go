package main

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os"

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

// Helper function to parse boolean values
func parseBool(value interface{}) bool {
    switch v := value.(type) {
    case bool:
        return v
    case string:
        return v == "true"
    default:
        panic(fmt.Sprintf("unexpected type for boolean value: %T", v))
    }
}

// Helper function to parse integer values
func parseInt(value interface{}) int {
    switch v := value.(type) {
    case float64:
        return int(v)
    case int:
        return v
    default:
        panic(fmt.Sprintf("unexpected type for integer value: %T", v))
    }
}

func handleLan(subdocData map[string]interface{}) EmbeddedLan {
    return EmbeddedLan{
        Lan: EmbeddedLanData{
            DhcpServerEnable:   parseBool(subdocData["DhcpServerEnable"]),
            LanIPAddress:       subdocData["LanIPAddress"].(string),
            LanSubnetMask:      subdocData["LanSubnetMask"].(string),
            DhcpStartIPAddress: subdocData["DhcpStartIPAddress"].(string),
            DhcpEndIPAddress:   subdocData["DhcpEndIPAddress"].(string),
            LeaseTime:          parseInt(subdocData["LeaseTime"]),
        },
    }
}

func handleWan(subdocData map[string]interface{}) EmbeddedWan {
    return EmbeddedWan{
        Wan: EmbeddedWanData{
            Enable:     parseBool(subdocData["Enable"]),
            InternalIP: subdocData["InternalIP"].(string),
        },
    }
}

func handlePortForwarding(subdocData interface{}) PortForwarding {
    var portForwardingList []PortForwardingData
    for _, item := range subdocData.([]interface{}) {
        portForwardingList = append(portForwardingList, PortForwardingData{
            InternalClient:       item.(map[string]interface{})["InternalClient"].(string),
            ExternalPortEndRange: item.(map[string]interface{})["ExternalPortEndRange"].(string),
            Enable:               parseBool(item.(map[string]interface{})["Enable"]),
            Protocol:             item.(map[string]interface{})["Protocol"].(string),
            Description:          item.(map[string]interface{})["Description"].(string),
            ExternalPort:         item.(map[string]interface{})["ExternalPort"].(string),
        })
    }
    return PortForwarding{
        PortForwarding: portForwardingList,
    }
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

    // Read JSON file
    data, err := ioutil.ReadFile(jsonFile)
    if err != nil {
        log.Fatalf("Failed to read JSON file: %v", err)
        return err
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
        subdoc = handleLan(subdocData)
    case "wan":
        subdoc = handleWan(subdocData)
    case "portforwarding":
        subdoc = handlePortForwarding(subdocData["listData"])
    case "privatessid":
        subdoc = subdocData
    default:
        subdoc = subdocData
    }

    prettyJSON, err := json.MarshalIndent(subdoc, "", "  ")
    fmt.Println(string(prettyJSON))

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

func writeToDB(macAddr string, subdocName string) error {
    // Prepare URL and headers, Replace with actual server host ip and port for webconfig server
    url := fmt.Sprintf("http://127.0.0.1:5000/api/v1/device/%s/document/%s", macAddr, subdocName)
    headers := map[string]string{
        "Content-type": "application/msgpack",
    }

    // Read MsgPack file
    msgpackBin := fmt.Sprintf("%s.msgpack", subdocName)
    fileData, err := ioutil.ReadFile(msgpackBin)
    if err != nil {
        return fmt.Errorf("failed to read MsgPack file: %v", err)
    }

    // Send POST request
    req, err := http.NewRequest("POST", url, ioutil.NopCloser(bytes.NewReader(fileData)))
    if err != nil {
        return fmt.Errorf("failed to create POST request: %v", err)
    }
    for key, value := range headers {
        req.Header.Set(key, value)
    }

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

    // Create MsgPack binary from JSON data
    err := createSubdoc(subdocJSON, subdocName, tr181)
    if err != nil {
        log.Fatalf("Error creating MsgPack binary: %v", err)
    }

    // Push MsgPack data to the server
    if err := writeToDB(macAddr, subdocName); err != nil {
        log.Fatalf("Error pushing MsgPack data: %v", err)
    }
}

