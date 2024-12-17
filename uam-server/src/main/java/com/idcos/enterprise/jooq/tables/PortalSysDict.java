/**
 * This class is generated by jOOQ
 */
package com.idcos.enterprise.jooq.tables;


import com.idcos.enterprise.jooq.Clouduam;
import com.idcos.enterprise.jooq.Keys;
import com.idcos.enterprise.jooq.tables.records.PortalSysDictRecord;

import java.util.Arrays;
import java.util.List;

import javax.annotation.Generated;

import org.jooq.Field;
import org.jooq.Table;
import org.jooq.TableField;
import org.jooq.UniqueKey;
import org.jooq.impl.TableImpl;


/**
 * uam系统参数表
 */
@Generated(
	value = {
		"http://www.jooq.org",
		"jOOQ version:3.7.2"
	},
	comments = "This class is generated by jOOQ"
)
@SuppressWarnings({ "all", "unchecked", "rawtypes" })
public class PortalSysDict extends TableImpl<PortalSysDictRecord> {

	private static final long serialVersionUID = -990219296;

	/**
	 * The reference instance of <code>clouduam.PORTAL_SYS_DICT</code>
	 */
	public static final PortalSysDict PORTAL_SYS_DICT = new PortalSysDict();

	/**
	 * The class holding records for this type
	 */
	@Override
	public Class<PortalSysDictRecord> getRecordType() {
		return PortalSysDictRecord.class;
	}

	/**
	 * The column <code>clouduam.PORTAL_SYS_DICT.TYPE_CODE</code>. 系统字典类型编码
	 */
	public final TableField<PortalSysDictRecord, String> TYPE_CODE = createField("TYPE_CODE", org.jooq.impl.SQLDataType.VARCHAR.length(64).nullable(false), this, "系统字典类型编码");

	/**
	 * The column <code>clouduam.PORTAL_SYS_DICT.CODE</code>. 系统字典编码
	 */
	public final TableField<PortalSysDictRecord, String> CODE = createField("CODE", org.jooq.impl.SQLDataType.VARCHAR.length(64).nullable(false), this, "系统字典编码");

	/**
	 * The column <code>clouduam.PORTAL_SYS_DICT.VALUE</code>. 参数值
	 */
	public final TableField<PortalSysDictRecord, String> VALUE = createField("VALUE", org.jooq.impl.SQLDataType.CLOB, this, "参数值");

	/**
	 * The column <code>clouduam.PORTAL_SYS_DICT.TENANT_ID</code>. 租户code
	 */
	public final TableField<PortalSysDictRecord, String> TENANT_ID = createField("TENANT_ID", org.jooq.impl.SQLDataType.VARCHAR.length(64).nullable(false), this, "租户code");

	/**
	 * The column <code>clouduam.PORTAL_SYS_DICT.REMARK</code>. 说明
	 */
	public final TableField<PortalSysDictRecord, String> REMARK = createField("REMARK", org.jooq.impl.SQLDataType.CLOB, this, "说明");

	/**
	 * Create a <code>clouduam.PORTAL_SYS_DICT</code> table reference
	 */
	public PortalSysDict() {
		this("PORTAL_SYS_DICT", null);
	}

	/**
	 * Create an aliased <code>clouduam.PORTAL_SYS_DICT</code> table reference
	 */
	public PortalSysDict(String alias) {
		this(alias, PORTAL_SYS_DICT);
	}

	private PortalSysDict(String alias, Table<PortalSysDictRecord> aliased) {
		this(alias, aliased, null);
	}

	private PortalSysDict(String alias, Table<PortalSysDictRecord> aliased, Field<?>[] parameters) {
		super(alias, Clouduam.CLOUDUAM, aliased, parameters, "uam系统参数表");
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public UniqueKey<PortalSysDictRecord> getPrimaryKey() {
		return Keys.KEY_PORTAL_SYS_DICT_PRIMARY;
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public List<UniqueKey<PortalSysDictRecord>> getKeys() {
		return Arrays.<UniqueKey<PortalSysDictRecord>>asList(Keys.KEY_PORTAL_SYS_DICT_PRIMARY);
	}

	/**
	 * {@inheritDoc}
	 */
	@Override
	public PortalSysDict as(String alias) {
		return new PortalSysDict(alias, this);
	}

	/**
	 * Rename this table
	 */
	public PortalSysDict rename(String name) {
		return new PortalSysDict(name, null);
	}
}
