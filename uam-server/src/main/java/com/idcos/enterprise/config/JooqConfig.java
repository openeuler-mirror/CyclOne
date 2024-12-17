package com.idcos.enterprise.config;

import org.jooq.DSLContext;
import org.jooq.SQLDialect;
import org.jooq.conf.Settings;
import org.jooq.conf.StatementType;
import org.jooq.impl.DSL;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import javax.sql.DataSource;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月03 上午9:51 souakiragen Exp $
 */

@Configuration
public class JooqConfig {

    @Bean
    public DSLContext dslContext(DataSource dataSource) {
        Settings settings = new Settings();
        settings.setStatementType(StatementType.PREPARED_STATEMENT);
        settings.setRenderSchema(false);
        settings.setRenderScalarSubqueriesForStoredFunctions(true);
        return DSL.using(dataSource, SQLDialect.MYSQL, settings);
    }

}
