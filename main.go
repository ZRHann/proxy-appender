package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"gopkg.in/yaml.v3"
)

func main() {
	http.HandleFunc("/xxxxxxxx/clash", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
			return
		}

		// 下载 YAML 配置文件
		client := http.DefaultClient
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
			return
		}
		req.Header.Set("User-Agent", "clash")
		
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to download file: %v", err), http.StatusInternalServerError)
			return
		}
		if resp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("Failed to download file: %v", resp.StatusCode), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		yamlData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read response body: %v", err), http.StatusInternalServerError)
			return
		}
		if len(yamlData) == 0 {
			http.Error(w, "Empty response body", http.StatusInternalServerError)
			return
		}

		// 解析 YAML 配置为 map[string]interface{} 以保留所有字段
		var config map[string]interface{}
		err = yaml.Unmarshal(yamlData, &config)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse YAML: %v", err), http.StatusInternalServerError)
			return
		}

		// 修改 proxies 字段，添加新的代理配置
		newProxy := map[string]interface{}{
			"type":     "vmess",
			"name":     "SG_AZURE",
			"server":   "xxxxxxxx",
			"port":     "xxxxxxxx",
			"uuid":     "xxxxxxxx",
			"alterId":  "0",
			"cipher":   "auto",
			"network":  "tcp",
		}
		
		// 获取现有的 proxies 列表
		if proxies, exists := config["proxies"]; exists {
			if proxiesList, ok := proxies.([]interface{}); ok {
				// 将 map 转换为 interface{} 并添加到列表
				proxiesList = append(proxiesList, newProxy)
				config["proxies"] = proxiesList
			}
		} else {
			config["proxies"] = []interface{}{newProxy}
		}

		// 修改 domain 字段，添加新的域名配置
		rulesToAppend := []interface{}{
			"DOMAIN,browser-intake-datadoghq.com,SG_AZURE",
			"DOMAIN,chat.openai.com.cdn.cloudflare.net,SG_AZURE",
			"DOMAIN,openai-api.arkoselabs.com,SG_AZURE",
			"DOMAIN,openaicom-api-bdcpf8c6d2e9atf6.z01.azurefd.net,SG_AZURE",
			"DOMAIN,openaicomproductionae4b.blob.core.windows.net,SG_AZURE",
			"DOMAIN,production-openaicom-storage.azureedge.net,SG_AZURE",
			"DOMAIN,static.cloudflareinsights.com,SG_AZURE",
			"DOMAIN-SUFFIX,ai.com,SG_AZURE",
			"DOMAIN-SUFFIX,algolia.net,SG_AZURE",
			"DOMAIN-SUFFIX,api.statsig.com,SG_AZURE",
			"DOMAIN-SUFFIX,auth0.com,SG_AZURE",
			"DOMAIN-SUFFIX,chatgpt.com,SG_AZURE",
			"DOMAIN-SUFFIX,chatgpt.livekit.cloud,SG_AZURE",
			"DOMAIN-SUFFIX,client-api.arkoselabs.com,SG_AZURE",
			"DOMAIN-SUFFIX,events.statsigapi.net,SG_AZURE",
			"DOMAIN-SUFFIX,featuregates.org,SG_AZURE",
			"DOMAIN-SUFFIX,host.livekit.cloud,SG_AZURE",
			"DOMAIN-SUFFIX,identrust.com,SG_AZURE",
			"DOMAIN-SUFFIX,intercom.io,SG_AZURE",
			"DOMAIN-SUFFIX,intercomcdn.com,SG_AZURE",
			"DOMAIN-SUFFIX,launchdarkly.com,SG_AZURE",
			"DOMAIN-SUFFIX,oaiusercontent.com,SG_AZURE",
			"DOMAIN-SUFFIX,observeit.net,SG_AZURE",
			"DOMAIN-SUFFIX,openai.com,SG_AZURE",
			"DOMAIN-SUFFIX,openaiapi-site.azureedge.net,SG_AZURE",
			"DOMAIN-SUFFIX,openaicom.imgix.net,SG_AZURE",
			"DOMAIN-SUFFIX,segment.io,SG_AZURE",
			"DOMAIN-SUFFIX,sentry.io,SG_AZURE",
			"DOMAIN-SUFFIX,stripe.com,SG_AZURE",
			"DOMAIN-SUFFIX,turn.livekit.cloud,SG_AZURE",
		}
		
		// 获取现有的 rules 列表
		if rules, exists := config["rules"]; exists {
			if rulesList, ok := rules.([]interface{}); ok {
				// 将新域名添加到列表开头
				config["rules"] = append(rulesToAppend, rulesList...)
			}
		} else {
			config["rules"] = rulesToAppend
		}

		// 序列化修改后的配置
		updatedYaml, err := yaml.Marshal(config)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to serialize updated YAML: %v", err), http.StatusInternalServerError)
			return
		}

		// 返回更新后的 YAML 配置
		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
		w.Write(updatedYaml)
	})

	// 启动 HTTP 服务
	port := "xxxxxxxx"
	fmt.Printf("Server started on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
