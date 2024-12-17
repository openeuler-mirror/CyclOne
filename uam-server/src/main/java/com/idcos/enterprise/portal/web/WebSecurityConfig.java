
package com.idcos.enterprise.portal.web;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.method.configuration.EnableGlobalMethodSecurity;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.web.authentication.AnonymousAuthenticationFilter;
import org.springframework.web.servlet.config.annotation.CorsRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurerAdapter;

import com.idcos.common.web.filter.AuthorizationFilter;
import com.idcos.common.web.service.AuthorizationProvider;

/**
 * @author zhouqin
 * @version com.idcos.enterprise.portal.web.security.WebSecurityConfig.java, v 1.1 12/29/15 zhouqin Exp $
 */
@Configuration
@EnableWebSecurity
@EnableGlobalMethodSecurity(securedEnabled = true)
public class WebSecurityConfig extends WebSecurityConfigurerAdapter {

    @Autowired
    private GlobalValue globalValue;

    @Autowired
    private AuthorizationProvider authorizationProvider;

    @Override
    protected void configure(HttpSecurity http) throws Exception {
        http.csrf().disable().sessionManagement().sessionCreationPolicy(SessionCreationPolicy.STATELESS);

        JwtFilter jwtFilter = new JwtFilter(globalValue.getSsoLoginUrl(), globalValue.getSecretKey());
        jwtFilter.setIgnoreUris(globalValue.getIgnoreUris());
        http.addFilterAfter(jwtFilter, AnonymousAuthenticationFilter.class);
        http.addFilterAfter(new AuthorizationFilter(authorizationProvider), JwtFilter.class);

    }

    @Bean
    public WebMvcConfigurer corsConfigurer() {
        return new WebMvcConfigurerAdapter() {
            @Override
            public void addCorsMappings(CorsRegistry registry) {
                registry.addMapping("/**");
            }
        };
    }

}
