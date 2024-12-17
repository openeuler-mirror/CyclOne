/**
 * This class is generated by jOOQ
 */
package com.idcos.enterprise.jooq.tables.records;


import com.idcos.enterprise.jooq.tables.PortalResource;

import javax.annotation.Generated;

import org.jooq.Field;
import org.jooq.Record1;
import org.jooq.Record8;
import org.jooq.Row8;
import org.jooq.impl.UpdatableRecordImpl;


/**
 * 权限资源表
 */
@Generated(
	value = {
		"http://www.jooq.org",
		"jOOQ version:3.7.2"
	},
	comments = "This class is generated by jOOQ"
)
@SuppressWarnings({ "all", "unchecked", "rawtypes" })
public class PortalResourceRecord extends UpdatableRecordImpl<PortalResourceRecord> implements Record8<String, String, String, String, String, String, String, String> {

	private static final long serialVersionUID = -581895275;

	/**
	 * Setter for <code>clouduam.PORTAL_RESOURCE.ID</code>.
	 */
	public void setId(String value) {
		setValue(0, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_RESOURCE.ID</code>.
	 */
	public String getId() {
		return (String) getValue(0);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_RESOURCE.APP_ID</code>. 应用系统名称
	 */
	public void setAppId(String value) {
		setValue(1, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_RESOURCE.APP_ID</code>. 应用系统名称
	 */
	public String getAppId() {
		return (String) getValue(1);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_RESOURCE.CODE</code>. 权限资源类型
	 */
	public void setCode(String value) {
		setValue(2, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_RESOURCE.CODE</code>. 权限资源类型
	 */
	public String getCode() {
		return (String) getValue(2);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_RESOURCE.NAME</code>. 权限资源名称
	 */
	public void setName(String value) {
		setValue(3, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_RESOURCE.NAME</code>. 权限资源名称
	 */
	public String getName() {
		return (String) getValue(3);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_RESOURCE.URL</code>. 权限资源URL
	 */
	public void setUrl(String value) {
		setValue(4, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_RESOURCE.URL</code>. 权限资源URL
	 */
	public String getUrl() {
		return (String) getValue(4);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_RESOURCE.REMARK</code>. 备注
	 */
	public void setRemark(String value) {
		setValue(5, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_RESOURCE.REMARK</code>. 备注
	 */
	public String getRemark() {
		return (String) getValue(5);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_RESOURCE.IS_ACTIVE</code>.
	 */
	public void setIsActive(String value) {
		setValue(6, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_RESOURCE.IS_ACTIVE</code>.
	 */
	public String getIsActive() {
		return (String) getValue(6);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_RESOURCE.TENANT</code>. 租户code
	 */
	public void setTenant(String value) {
		setValue(7, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_RESOURCE.TENANT</code>. 租户code
	 */
	public String getTenant() {
		return (String) getValue(7);
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
	// Record8 type implementation
	// -------------------------------------------------------------------------

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Row8<String, String, String, String, String, String, String, String> fieldsRow() {
		return (Row8) super.fieldsRow();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Row8<String, String, String, String, String, String, String, String> valuesRow() {
		return (Row8) super.valuesRow();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field1() {
		return PortalResource.PORTAL_RESOURCE.ID;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field2() {
		return PortalResource.PORTAL_RESOURCE.APP_ID;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field3() {
		return PortalResource.PORTAL_RESOURCE.CODE;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field4() {
		return PortalResource.PORTAL_RESOURCE.NAME;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field5() {
		return PortalResource.PORTAL_RESOURCE.URL;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field6() {
		return PortalResource.PORTAL_RESOURCE.REMARK;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field7() {
		return PortalResource.PORTAL_RESOURCE.IS_ACTIVE;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field8() {
		return PortalResource.PORTAL_RESOURCE.TENANT;
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
		return getAppId();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value3() {
		return getCode();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value4() {
		return getName();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value5() {
		return getUrl();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value6() {
		return getRemark();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value7() {
		return getIsActive();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value8() {
		return getTenant();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalResourceRecord value1(String value) {
		setId(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalResourceRecord value2(String value) {
		setAppId(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalResourceRecord value3(String value) {
		setCode(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalResourceRecord value4(String value) {
		setName(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalResourceRecord value5(String value) {
		setUrl(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalResourceRecord value6(String value) {
		setRemark(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalResourceRecord value7(String value) {
		setIsActive(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalResourceRecord value8(String value) {
		setTenant(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalResourceRecord values(String value1, String value2, String value3, String value4, String value5, String value6, String value7, String value8) {
		value1(value1);
		value2(value2);
		value3(value3);
		value4(value4);
		value5(value5);
		value6(value6);
		value7(value7);
		value8(value8);
		return this;
	}

	// -------------------------------------------------------------------------
	// Constructors
	// -------------------------------------------------------------------------

	/**
	 * Create a detached PortalResourceRecord
	 */
	public PortalResourceRecord() {
		super(PortalResource.PORTAL_RESOURCE);
	}

	/**
	 * Create a detached, initialised PortalResourceRecord
	 */
	public PortalResourceRecord(String id, String appId, String code, String name, String url, String remark, String isActive, String tenant) {
		super(PortalResource.PORTAL_RESOURCE);

		setValue(0, id);
		setValue(1, appId);
		setValue(2, code);
		setValue(3, name);
		setValue(4, url);
		setValue(5, remark);
		setValue(6, isActive);
		setValue(7, tenant);
	}
}
