
package com.idcos.enterprise.portal.form;

import com.idcos.cloud.core.dal.common.page.PageForm;
import com.idcos.cloud.core.dal.common.query.DataQueryField;
import com.idcos.cloud.core.dal.common.query.OperatorEnum;

/**
 * 权限资源管理通用查询form
 *
 * @author pengganyu
 * @version $Id: PortalQueryByPageForm.java, v 0.1 2016年6月1日 下午7:27:36 pengganyu Exp $
 */
public class PortalQueryByPageForm extends PageForm {

    @DataQueryField
    private String id;

    @DataQueryField(operator = OperatorEnum.LIKE)
    private String cnd;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getCnd() {
        return cnd;
    }

    public void setCnd(String cnd) {
        this.cnd = cnd;
    }

}
