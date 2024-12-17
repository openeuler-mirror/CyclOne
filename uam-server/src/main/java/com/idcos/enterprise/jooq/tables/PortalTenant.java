/**
 * This class is generated by jOOQ
 */
package com.idcos.enterprise.jooq.tables;


import com.idcos.enterprise.jooq.Clouduam;
import com.idcos.enterprise.jooq.Keys;
import com.idcos.enterprise.jooq.tables.records.PortalTenantRecord;

import java.sql.Timestamp;
import java.util.Arrays;
import java.util.List;

import javax.annotation.Generated;

import org.jooq.Field;
import org.jooq.Table;
import org.jooq.TableField;
import org.jooq.UniqueKey;
import org.jooq.impl.TableImpl;


/**
 * 租户信息表
 */
@Generated(
	value = {
		"http://www.jooq.org",
		"jOOQ version:3.7.2"
	},
	comments = "This class is generated by jOOQ"
)
@SuppressWarnings({ "all", "unchecked", "rawtypes" })
public class PortalTenant extends TableImpl<PortalTenantRecord> {

	private static final long serialVersionUID = -892349280;

	/**
	 * The reference instance of <code>clouduam.PORTAL_TENANT</code>
	 */
	public static final PortalTenant PORTAL_TENANT = new PortalTenant();

	/**
	 * The class holding records for this type
	 */
	@Override
	public Class<PortalTenantRecord> getRecordType() {
		return PortalTenantRecord.class;
	}

	/**
	 * The column <code>clouduam.PORTAL_TENANT.ID</code>.
	 */
	public final TableField<PortalTenantRecord, String> ID = createField("ID", org.jooq.impl.SQLDataType.VARCHAR.length(64).nullable(false), this, "");

	/**
	 * The column <code>clouduam.PORTAL_TENANT.NAME</code>. 租户编码
	 */
	public final TableField<PortalTenantRecord, String> NAME = createField("NAME", org.jooq.impl.SQLDataType.VARCHAR.length(128), this, "租户编码");

	/**
	 * The column <code>clouduam.PORTAL_TENANT.DISPLAY_NAME</code>. 租户名称
	 */
	public final TableField<PortalTenantRecord, String> DISPLAY_NAME = createField("DISPLAY_NAME", org.jooq.impl.SQLDataType.VARCHAR.length(128), this, "租户名称");

	/**
	 * The column <code>clouduam.PORTAL_TENANT.GMT_CREATE</code>. 创建日期
	 */
	public final TableField<PortalTenantRecord, Timestamp> GMT_CREATE = createField("GMT_CREATE", org.jooq.impl.SQLDataType.TIMESTAMP, this, "创建日期");

	/**
	 * The column <code>clouduam.PORTAL_TENANT.GMT_MODIFIED</code>. 修改日期
	 */
	public final TableField<PortalTenantRecord, Timestamp> GMT_MODIFIED = createField("GMT_MODIFIED", org.jooq.impl.SQLDataType.TIMESTAMP.nullable(false), this, "修改日期");

	/**
	 * The column <code>clouduam.PORTAL_TENANT.IS_ACTIVE</code>. 是否可用
	 */
	public final TableField<PortalTenantRecord, String> IS_ACTIVE = createField("IS_ACTIVE", org.jooq.impl.SQLDataType.CHAR.length(1), this, "是否可用");

	/**
	 * Create a <code>clouduam.PORTAL_TENANT</code> table reference
	 */
	public PortalTenant() {
		this("PORTAL_TENANT", null);
	}

	/**
	 * Create an aliased <code>clouduam.PORTAL_TENANT</code> table reference
	 */
	public PortalTenant(String alias) {
		this(alias, PORTAL_TENANT);
	}

	private PortalTenant(String alias, Table<PortalTenantRecord> aliased) {
		this(alias, aliased, null);
	}

	private PortalTenant(String alias, Table<PortalTenantRecord> aliased, Field<?>[] parameters) {
		super(alias, Clouduam.CLOUDUAM, aliased, parameters, "租户信息表");
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public UniqueKey<PortalTenantRecord> getPrimaryKey() {
		return Keys.KEY_PORTAL_TENANT_PRIMARY;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public List<UniqueKey<PortalTenantRecord>> getKeys() {
		return Arrays.<UniqueKey<PortalTenantRecord>>asList(Keys.KEY_PORTAL_TENANT_PRIMARY);
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalTenant as(String alias) {
		return new PortalTenant(alias, this);
	}

	/**
	 * Rename this table
	 */
	public PortalTenant rename(String name) {
		return new PortalTenant(name, null);
	}
}