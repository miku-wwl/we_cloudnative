package com.weilai.callservice;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

@RestController
@SpringBootApplication
public class CallServiceApplication {

    private final RestTemplate restTemplate = new RestTemplate();

    public static void main(String[] args) {
        SpringApplication.run(CallServiceApplication.class, args);
    }

    @GetMapping("/call-target")
    public String callTarget() {
        String url = "http://target-service/hello";
        return restTemplate.getForObject(url, String.class);
    }

    @GetMapping("/hello")
    public String hello() {
        return "hello world!";
    }

}
