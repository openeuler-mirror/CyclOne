/**
 * This class is generated by jOOQ
 */
package com.idcos.enterprise.jooq.tables;


import com.idcos.enterprise.jooq.Clouduam;
import com.idcos.enterprise.jooq.Keys;
import com.idcos.enterprise.jooq.tables.records.PortalResourceRecord;

import java.util.Arrays;
import java.util.List;

import javax.annotation.Generated;

import org.jooq.Field;
import org.jooq.Table;
import org.jooq.TableField;
import org.jooq.UniqueKey;
import org.jooq.impl.TableImpl;


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
public class PortalResource extends TableImpl<PortalResourceRecord> {

	private static final long serialVersionUID = -243160886;

	/**
	 * The reference instance of <code>clouduam.PORTAL_RESOURCE</code>
	 */
	public static final PortalResource PORTAL_RESOURCE = new PortalResource();

	/**
	 * The class holding records for this type
	 */
	@Override
	public Class<PortalResourceRecord> getRecordType() {
		return PortalResourceRecord.class;
	}

	/**
	 * The column <code>clouduam.PORTAL_RESOURCE.ID</code>.
	 */
	public final TableField<PortalResourceRecord, String> ID = createField("ID", org.jooq.impl.SQLDataType.VARCHAR.length(64).nullable(false), this, "");

	/**
	 * The column <code>clouduam.PORTAL_RESOURCE.APP_ID</code>. 应用系统名称
	 */
	public final TableField<PortalResourceRecord, String> APP_ID = createField("APP_ID", org.jooq.impl.SQLDataType.VARCHAR.length(64).nullable(false), this, "应用系统名称");

	/**
	 * The column <code>clouduam.PORTAL_RESOURCE.CODE</code>. 权限资源类型
	 */
	public final TableField<PortalResourceRecord, String> CODE = createField("CODE", org.jooq.impl.SQLDataType.VARCHAR.length(64).nullable(false), this, "权限资源类型");

	/**
	 * The column <code>clouduam.PORTAL_RESOURCE.NAME</code>. 权限资源名称
	 */
	public final TableField<PortalResourceRecord, String> NAME = createField("NAME", org.jooq.impl.SQLDataType.VARCHAR.length(128).nullable(false), this, "权限资源名称");

	/**
	 * The column <code>clouduam.PORTAL_RESOURCE.URL</code>. 权限资源URL
	 */
	public final TableField<PortalResourceRecord, String> URL = createField("URL", org.jooq.impl.SQLDataType.VARCHAR.length(128).nullable(false), this, "权限资源URL");

	/**
	 * The column <code>clouduam.PORTAL_RESOURCE.REMARK</code>. 备注
	 */
	public final TableField<PortalResourceRecord, String> REMARK = createField("REMARK", org.jooq.impl.SQLDataType.VARCHAR.length(256).nullable(false), this, "备注");

	/**
	 * The column <code>clouduam.PORTAL_RESOURCE.IS_ACTIVE</code>.
	 */
	public final TableField<PortalResourceRecord, String> IS_ACTIVE = createField("IS_ACTIVE", org.jooq.impl.SQLDataType.VARCHAR.length(1).nullable(false).defaulted(true), this, "");

	/**
	 * The column <code>clouduam.PORTAL_RESOURCE.TENANT</code>. 租户code
	 */
	public final TableField<PortalResourceRecord, String> TENANT = createField("TENANT", org.jooq.impl.SQLDataType.VARCHAR.length(64), this, "租户code");

	/**
	 * Create a <code>clouduam.PORTAL_RESOURCE</code> table reference
	 */
	public PortalResource() {
		this("PORTAL_RESOURCE", null);
	}

	/**
	 * Create an aliased <code>clouduam.PORTAL_RESOURCE</code> table reference
	 */
	public PortalResource(String alias) {
		this(alias, PORTAL_RESOURCE);
	}

	private PortalResource(String alias, Table<PortalResourceRecord> aliased) {
		this(alias, aliased, null);
	}

	private PortalResource(String alias, Table<PortalResourceRecord> aliased, Field<?>[] parameters) {
		super(alias, Clouduam.CLOUDUAM, aliased, parameters, "权限资源表");
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public UniqueKey<PortalResourceRecord> getPrimaryKey() {
		return Keys.KEY_PORTAL_RESOURCE_PRIMARY;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public List<UniqueKey<PortalResourceRecord>> getKeys() {
		return Arrays.<UniqueKey<PortalResourceRecord>>asList(Keys.KEY_PORTAL_RESOURCE_PRIMARY);
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalResource as(String alias) {
		return new PortalResource(alias, this);
	}

	/**
	 * Rename this table
	 */
	public PortalResource rename(String name) {
		return new PortalResource(name, null);
	}
}
