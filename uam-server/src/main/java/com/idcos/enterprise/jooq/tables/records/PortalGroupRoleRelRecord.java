/**
 * This class is generated by jOOQ
 */
package com.idcos.enterprise.jooq.tables.records;


import com.idcos.enterprise.jooq.tables.PortalGroupRoleRel;

import javax.annotation.Generated;

import org.jooq.Field;
import org.jooq.Record1;
import org.jooq.Record4;
import org.jooq.Row4;
import org.jooq.impl.UpdatableRecordImpl;


/**
 * 用户组与角色关系表
 */
@Generated(
	value = {
		"http://www.jooq.org",
		"jOOQ version:3.7.2"
	},
	comments = "This class is generated by jOOQ"
)
@SuppressWarnings({ "all", "unchecked", "rawtypes" })
public class PortalGroupRoleRelRecord extends UpdatableRecordImpl<PortalGroupRoleRelRecord> implements Record4<String, String, String, String> {

	private static final long serialVersionUID = -1466261113;

	/**
	 * Setter for <code>clouduam.PORTAL_GROUP_ROLE_REL.ID</code>. 关系ID
	 */
	public void setId(String value) {
		setValue(0, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_GROUP_ROLE_REL.ID</code>. 关系ID
	 */
	public String getId() {
		return (String) getValue(0);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_GROUP_ROLE_REL.ROLE_ID</code>. 角色ID
	 */
	public void setRoleId(String value) {
		setValue(1, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_GROUP_ROLE_REL.ROLE_ID</code>. 角色ID
	 */
	public String getRoleId() {
		return (String) getValue(1);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_GROUP_ROLE_REL.GROUP_ID</code>. 用户组ID
	 */
	public void setGroupId(String value) {
		setValue(2, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_GROUP_ROLE_REL.GROUP_ID</code>. 用户组ID
	 */
	public String getGroupId() {
		return (String) getValue(2);
	}

	/**
	 * Setter for <code>clouduam.PORTAL_GROUP_ROLE_REL.TENANT</code>. 租户code
	 */
	public void setTenant(String value) {
		setValue(3, value);
	}

	/**
	 * Getter for <code>clouduam.PORTAL_GROUP_ROLE_REL.TENANT</code>. 租户code
	 */
	public String getTenant() {
		return (String) getValue(3);
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
	// Record4 type implementation
	// -------------------------------------------------------------------------

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Row4<String, String, String, String> fieldsRow() {
		return (Row4) super.fieldsRow();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Row4<String, String, String, String> valuesRow() {
		return (Row4) super.valuesRow();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field1() {
		return PortalGroupRoleRel.PORTAL_GROUP_ROLE_REL.ID;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field2() {
		return PortalGroupRoleRel.PORTAL_GROUP_ROLE_REL.ROLE_ID;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field3() {
		return PortalGroupRoleRel.PORTAL_GROUP_ROLE_REL.GROUP_ID;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public Field<String> field4() {
		return PortalGroupRoleRel.PORTAL_GROUP_ROLE_REL.TENANT;
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
		return getRoleId();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value3() {
		return getGroupId();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public String value4() {
		return getTenant();
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalGroupRoleRelRecord value1(String value) {
		setId(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalGroupRoleRelRecord value2(String value) {
		setRoleId(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalGroupRoleRelRecord value3(String value) {
		setGroupId(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalGroupRoleRelRecord value4(String value) {
		setTenant(value);
		return this;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalGroupRoleRelRecord values(String value1, String value2, String value3, String value4) {
		value1(value1);
		value2(value2);
		value3(value3);
		value4(value4);
		return this;
	}

	// -------------------------------------------------------------------------
	// Constructors
	// -------------------------------------------------------------------------

	/**
	 * Create a detached PortalGroupRoleRelRecord
	 */
	public PortalGroupRoleRelRecord() {
		super(PortalGroupRoleRel.PORTAL_GROUP_ROLE_REL);
	}

	/**
	 * Create a detached, initialised PortalGroupRoleRelRecord
	 */
	public PortalGroupRoleRelRecord(String id, String roleId, String groupId, String tenant) {
		super(PortalGroupRoleRel.PORTAL_GROUP_ROLE_REL);

		setValue(0, id);
		setValue(1, roleId);
		setValue(2, groupId);
		setValue(3, tenant);
	}
}
