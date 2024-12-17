/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web.auto.excelfilehandler;

import org.apache.poi.hssf.usermodel.HSSFDateUtil;
import org.apache.poi.ss.usermodel.CellType;
import org.apache.poi.ss.usermodel.Cell;

import java.text.SimpleDateFormat;

/**
 * 操作EXCEL的工具类
 *
 * @author Dana
 * @version ExcelUtil.java, v1 2017/12/10 下午6:10 Dana Exp $$
 */
public class ExcelUtil {

    /**
     *
     */
    private ExcelUtil() {
    }

    public static final String getValueFromCell(Cell cell) {
        if (cell == null) {
            return null;
        }
        String value;
        switch (cell.getCellType()) {
            case NUMERIC:
                // 数字
                if (HSSFDateUtil.isCellDateFormatted(cell)) {
                    // 如果是日期类型
                    value = new SimpleDateFormat("yyyy-MM-dd hh:mm:ss").format(cell.getDateCellValue());
                } else {
                    cell.setCellType(CellType.STRING);
                    value = cell.getStringCellValue();
                }
                break;
            case STRING:
                // 字符串
                value = cell.getStringCellValue();
                break;
            case FORMULA:
                // 用数字方式获取公式结果，根据值判断是否为日期类型
                double numericValue = cell.getNumericCellValue();
                if (HSSFDateUtil.isValidExcelDate(numericValue)) {
                    // 如果是日期类型
                    value = new SimpleDateFormat("yyyy-MM-dd hh:mm:ss").format(cell.getDateCellValue());
                } else {
                    value = String.valueOf(numericValue);
                }
                break;
            case BLANK:
                // 空白
                value = "";
                break;
            case BOOLEAN:
                // Boolean
                value = String.valueOf(cell.getBooleanCellValue());
                break;
            case ERROR:
                // Error，返回错误码
                value = String.valueOf(cell.getErrorCellValue());
                break;
            default:
                value = "";
                break;
        }
        // 使用[]记录坐标
        return value;
    }
}