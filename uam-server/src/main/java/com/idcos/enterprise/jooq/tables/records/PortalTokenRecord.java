/**
 * This class is generated by jOOQ
 */
package com.idcos.enterprise.jooq.tables.records;


import com.idcos.enterprise.jooq.tables.PortalToken;

import java.sql.Timestamp;

import javax.annotation.Generated;

import org.jooq.Field;
import org.jooq.Record1;
import org.jooq.Record10;
import org.jooq.Row10;
import org.jooq.impl.UpdatableRecordImpl;
import org.jooq.types.UInteger;


/**
 * token信息表
 */
@Generated(
	value = {
		"http://www.jooq.org",
		"jOOQ version:3.7.2"
	},
	comments = "This class is generated by jOOQ"
)
@SuppressWarnings({ "all", "unchecked", "rawtypes" })
public class PortalTokenRecord extends UpdatableRecordImpl<PortalTokenRecord> implements Record10<String, String, UInteger, String, String, String, Timestamp, Timestamp, Timestamp, String> {

	private static final long serialVersionUID = 1273843376;

	/**
	 * Setter for <code>clouduam.PORTAL_TOKEN.ID</code>.
	 */
	public void setId(String value) {
		setValue(0, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_TOKEN.ID</code>.
	 */
	public String getId() {
		return (String) getValue(0);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_TOKEN.NAME</code>. token值
	 */
	public void setName(String value) {
		setValue(1, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_TOKEN.NAME</code>. token值
	 */
	public String getName() {
		return (String) getValue(1);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_TOKEN.TOKEN_CRC</code>. token串的crc哈希值
	 */
	public void setTokenCrc(UInteger value) {
		setValue(2, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_TOKEN.TOKEN_CRC</code>. token串的crc哈希值
	 */
	public UInteger getTokenCrc() {
		return (UInteger) getValue(2);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_TOKEN.LOGIN_ID</code>. 登录名
	 */
	public void setLoginId(String value) {
		setValue(3, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_TOKEN.LOGIN_ID</code>. 登录名
	 */
	public String getLoginId() {
		return (String) getValue(3);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_TOKEN.TENANT_ID</code>. 租户code
	 */
	public void setTenantId(String value) {
		setValue(4, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_TOKEN.TENANT_ID</code>. 租户code
	 */
	public String getTenantId() {
		return (String) getValue(4);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_TOKEN.IS_ACTIVE</code>. 是否可用
	 */
	public void setIsActive(String value) {
		setValue(5, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_TOKEN.IS_ACTIVE</code>. 是否可用
	 */
	public String getIsActive() {
		return (String) getValue(5);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_TOKEN.EXPIRE_TIME</code>. token过期时间
	 */
	public void setExpireTime(Timestamp value) {
		setValue(6, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_TOKEN.EXPIRE_TIME</code>. token过期时间
	 */
	public Timestamp getExpireTime() {
		return (Timestamp) getValue(6);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_TOKEN.GMT_CREATE</code>. 创建日期
	 */
	public void setGmtCreate(Timestamp value) {
		setValue(7, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_TOKEN.GMT_CREATE</code>. 创建日期
	 */
	public Timestamp getGmtCreate() {
		return (Timestamp) getValue(7);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_TOKEN.GMT_MODIFIED</code>. 修改日期
	 */
	public void setGmtModified(Timestamp value) {
		setValue(8, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_TOKEN.GMT_MODIFIED</code>. 修改日期
	 */
	public Timestamp getGmtModified() {
		return (Timestamp) getValue(8);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_TOKEN.REMARK</code>. 备注
	 */
	public void setRemark(String value) {
		setValue(9, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_TOKEN.REMARK</code>. 备注
	 */
	public String getRemark() {
		return (String) getValue(9);
	}

	// -------------------------------------------------------------------------
	// Primary key information
	// -------------------------------------------------------------------------

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Record1<String> key() {
		return (Record1) super.key();
	}

	// -------------------------------------------------------------------------
	// Record10 type implementation
	// -------------------------------------------------------------------------

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Row10<String, String, UInteger, String, String, String, Timestamp, Timestamp, Timestamp, String> fieldsRow() {
		return (Row10) super.fieldsRow();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Row10<String, String, UInteger, String, String, String, Timestamp, Timestamp, Timestamp, String> valuesRow() {
		return (Row10) super.valuesRow();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field1() {
		return PortalToken.PORTAL_TOKEN.ID;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field2() {
		return PortalToken.PORTAL_TOKEN.NAME;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<UInteger> field3() {
		return PortalToken.PORTAL_TOKEN.TOKEN_CRC;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field4() {
		return PortalToken.PORTAL_TOKEN.LOGIN_ID;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field5() {
		return PortalToken.PORTAL_TOKEN.TENANT_ID;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field6() {
		return PortalToken.PORTAL_TOKEN.IS_ACTIVE;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<Timestamp> field7() {
		return PortalToken.PORTAL_TOKEN.EXPIRE_TIME;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<Timestamp> field8() {
		return PortalToken.PORTAL_TOKEN.GMT_CREATE;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<Timestamp> field9() {
		return PortalToken.PORTAL_TOKEN.GMT_MODIFIED;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field10() {
		return PortalToken.PORTAL_TOKEN.REMARK;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value1() {
		return getId();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value2() {
		return getName();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public UInteger value3() {
		return getTokenCrc();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value4() {
		return getLoginId();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value5() {
		return getTenantId();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value6() {
		return getIsActive();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Timestamp value7() {
		return getExpireTime();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Timestamp value8() {
		return getGmtCreate();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Timestamp value9() {
		return getGmtModified();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value10() {
		return getRemark();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalTokenRecord value1(String value) {
		setId(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalTokenRecord value2(String value) {
		setName(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalTokenRecord value3(UInteger value) {
		setTokenCrc(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalTokenRecord value4(String value) {
		setLoginId(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalTokenRecord value5(String value) {
		setTenantId(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalTokenRecord value6(String value) {
		setIsActive(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalTokenRecord value7(Timestamp value) {
		setExpireTime(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalTokenRecord value8(Timestamp value) {
		setGmtCreate(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalTokenRecord value9(Timestamp value) {
		setGmtModified(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalTokenRecord value10(String value) {
		setRemark(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalTokenRecord values(String value1, String value2, UInteger value3, String value4, String value5, String value6, Timestamp value7, Timestamp value8, Timestamp value9, String value10) {
		value1(value1);
		value2(value2);
		value3(value3);
		value4(value4);
		value5(value5);
		value6(value6);
		value7(value7);
		value8(value8);
		value9(value9);
		value10(value10);
		return this;
	}

	// -------------------------------------------------------------------------
	// Constructors
	// -------------------------------------------------------------------------

	/**
	 * Create a detached PortalTokenRecord
	 */
	public PortalTokenRecord() {
		super(PortalToken.PORTAL_TOKEN);
	}

	/**
	 * Create a detached, initialised PortalTokenRecord
	 */
	public PortalTokenRecord(String id, String name, UInteger tokenCrc, String loginId, String tenantId, String isActive, Timestamp expireTime, Timestamp gmtCreate, Timestamp gmtModified, String remark) {
		super(PortalToken.PORTAL_TOKEN);

		setValue(0, id);
		setValue(1, name);
		setValue(2, tokenCrc);
		setValue(3, loginId);
		setValue(4, tenantId);
		setValue(5, isActive);
		setValue(6, expireTime);
		setValue(7, gmtCreate);
		setValue(8, gmtModified);
		setValue(9, remark);
	}
}
