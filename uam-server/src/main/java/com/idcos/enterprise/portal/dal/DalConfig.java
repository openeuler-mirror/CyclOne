package com.idcos.enterprise.portal.dal;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.jpa.repository.config.EnableJpaAuditing;
import org.springframework.transaction.PlatformTransactionManager;
import org.springframework.transaction.TransactionDefinition;
import org.springframework.transaction.support.TransactionTemplate;

/**
 * @author zhouqin
 * @version com.idcos.automate.dal.DalConfig.java, v 1.1 3/1/16 zhouqin Exp $
 */
@Configuration
@EnableJpaAuditing
public class DalConfig {

    @Bean(name = "transactionTemplateRequiredNew")
    public TransactionTemplate transactionTemplateRequiredNew(PlatformTransactionManager transactionManager) {
        TransactionTemplate transactionTemplate = new TransactionTemplate(transactionManager);
        transactionTemplate.setPropagationBehavior(TransactionDefinition.PROPAGATION_REQUIRES_NEW);

        return transactionTemplate;
    }
}
