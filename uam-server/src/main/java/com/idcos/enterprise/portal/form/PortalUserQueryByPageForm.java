
package com.idcos.enterprise.portal.form;

import com.idcos.cloud.biz.common.BaseForm;
import com.idcos.cloud.core.dal.common.query.DataQueryField;
import com.idcos.cloud.core.dal.common.query.OperatorEnum;

/**
 * 权限资源管理通用查询form
 *
 * @author pengganyu
 * @version $Id: PortalQueryByPageForm.java, v 0.1 2016年6月1日 下午7:27:36 pengganyu Exp $
 */
public class PortalUserQueryByPageForm extends BaseForm {

    @DataQueryField
    private String id;

    @DataQueryField(operator = OperatorEnum.LIKE)
    private String nameCn;

    @DataQueryField(operator = OperatorEnum.LIKE)
    private String nameEn;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getNameCn() {
        return nameCn;
    }

    public void setNameCn(String nameCn) {
        this.nameCn = nameCn;
    }

    public String getNameEn() {
        return nameEn;
    }

    public void setNameEn(String nameEn) {
        this.nameEn = nameEn;
    }

}
