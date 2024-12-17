/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web.auto.excelfilehandler;

import com.idcos.cloud.core.common.util.StringUtil;
import com.idcos.enterprise.portal.biz.common.utils.SourceTypeUtil;
import com.idcos.enterprise.portal.dal.entity.PortalDept;
import com.idcos.enterprise.portal.dal.entity.PortalUser;
import com.idcos.enterprise.portal.dal.repository.PortalDeptRepository;
import com.idcos.enterprise.portal.dal.repository.PortalUserRepository;
import org.apache.poi.hssf.usermodel.HSSFCellStyle;
import org.apache.poi.hssf.usermodel.HSSFFont;
import org.apache.poi.hssf.usermodel.HSSFWorkbook;
import org.apache.poi.ss.usermodel.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.sql.Timestamp;
import java.text.SimpleDateFormat;
import java.util.Arrays;
import java.util.List;

/**
 * @author Dana
 * @version UserExcelTempterHandler.java, v1 2017/12/12 下午11:01 Dana Exp $$
 */
@Service
public class UserExcelTempleteHandler {

    @Autowired
    private PortalUserRepository portalUserRepository;

    @Autowired
    private PortalDeptRepository portalDeptRepository;

    /**
     * 用户导入的excel表头模板
     */
    private static final String[] USER_COLUMNS = {"登录名", "姓名", "具体部门名称", "职务", "密码", "邮箱", "移动电话", "办公电话", "数据来源", "备注"};

    /**
     * 格式化时间到分钟。
     */
    private SimpleDateFormat sdf = new SimpleDateFormat("yyyy_MM_dd_HH_mm");

    /**
     * 下载用户信息的模板。
     *
     * @param request
     * @param response
     * @throws IOException
     */
    public void downloadTemplete(HttpServletRequest request, HttpServletResponse response) throws IOException {
        HSSFWorkbook wb = new HSSFWorkbook();
        Sheet sh = wb.createSheet("用户信息模板");
        //设置用户信息表的表头
        setHead(wb, sh);
        response.reset();
        // 设置response的Header
        response.setHeader("Pragma", "No-cache");
        response.setHeader("Cache-Control", "No-cache");
        response.setDateHeader("Expires", 0);
        response.setContentType("application/octet-stream");
        response.addHeader("Content-Disposition",
                "attachment;filename=" + new String(
                        ("用户信息模板_" + sdf.format(new Timestamp(System.currentTimeMillis())) + ".xls").getBytes("GBK"),
                        "iso8859-1"));
        response.setContentType("application/msexcel;charset=GBK");

        wb.write(response.getOutputStream());

        wb.close();
    }

    /**
     * 设置用户信息表的表头
     *
     * @param wb 代表一个EXCEL
     * @param sh 代表一个sheet页。
     */
    private void setHead(HSSFWorkbook wb, Sheet sh) {
        Row headRow = sh.createRow(0);

        headRow.setHeight((short) 400);

        HSSFFont rowF = wb.createFont();
        //字号
        rowF.setFontHeightInPoints((short) 10);
        //rowF.setBoldweight(HSSFFont.BOLDWEIGHT_NORMAL);
        rowF.setBold(true);
        //设置头部字体为宋体
        rowF.setFontName("宋体");

        HSSFCellStyle rowStyle = wb.createCellStyle();
        rowStyle.setFont(rowF);
        //左右居中
        rowStyle.setAlignment(HorizontalAlignment.CENTER);
        rowStyle.setVerticalAlignment(VerticalAlignment.CENTER);
        rowStyle.setFillBackgroundColor(IndexedColors.WHITE.getIndex());
        rowStyle.setWrapText(true);

        headRow.setRowStyle(rowStyle);

        HSSFCellStyle style = wb.createCellStyle();
        HSSFFont f = wb.createFont();
        //字号
        f.setFontHeightInPoints((short) 11);
        //加粗
        //f.setBoldweight(HSSFFont.BOLDWEIGHT_NORMAL);
        f.setBold(true);
        //设置头部字体为宋体
        f.setFontName("宋体");
        style.setFont(f);
        //左右居中
        style.setAlignment(HorizontalAlignment.CENTER);
        style.setVerticalAlignment(VerticalAlignment.CENTER);
        style.setFillBackgroundColor(IndexedColors.WHITE.getIndex());
        style.setWrapText(true);

        //必填项字体为红色

        //style.setBottomBorderColor(IndexedColors.BLACK.getIndex()); // 底部边框颜色

        HSSFCellStyle style1 = wb.createCellStyle();
        style1.setFont(f);
        //左右居中
        style1.setAlignment(HorizontalAlignment.CENTER);
        style1.setVerticalAlignment(VerticalAlignment.CENTER);
        HSSFFont font = wb.createFont();
        font.setColor(IndexedColors.RED.getIndex());
        //设置头部字体为宋体
        font.setFontName("宋体");
        //字号
        font.setFontHeightInPoints((short) 11);
        //加粗
        font.setBold(true);
        //font.setBoldweight(HSSFFont.BOLDWEIGHT_NORMAL);
        style1.setFont(font);
        style1.setFillBackgroundColor(IndexedColors.RED.getIndex());
        style1.setWrapText(true);

        CreationHelper factory = wb.getCreationHelper();
        // When the comment box is visible, have it show in a 1x3 space
        ClientAnchor anchor = factory.createClientAnchor();
        Drawing drawing = sh.createDrawingPatriarch();
        for (int i = 0; i < USER_COLUMNS.length; i++) {
            sh.setColumnWidth(i, 12 * 256);
            sh.autoSizeColumn(i);
            Cell cell = headRow.createCell(i);
            cell.setCellValue(USER_COLUMNS[i]);
            cell.setCellType(CellType.STRING);

            anchor.setCol1(cell.getColumnIndex());
            anchor.setCol2(cell.getColumnIndex() + 4);
            anchor.setRow1(cell.getRowIndex());
            anchor.setRow2(cell.getRowIndex() + 4);
            anchor.setDx1(0);
            anchor.setDx2(0);
            anchor.setDy1(1300);
            anchor.setDy2(1000);

            if (i < 2) {
                Comment comment = drawing.createCellComment(anchor);
                RichTextString str = factory.createRichTextString("不能为空。");
                comment.setString(str);
                comment.setAuthor("系统");
                str.applyFont(font);
                cell.setCellStyle(style1);
                cell.setCellComment(comment);
            } else {
                cell.setCellStyle(style);
            }
        }

    }

    /**
     * 下载用户数据。
     *
     * @param response
     * @throws IOException
     */
    public void downloadData(HttpServletResponse response, String deptId,
                             String tenantId) throws IOException {
        HSSFWorkbook wb = new HSSFWorkbook();
        Sheet sh = wb.createSheet("用户列表数据");
        //表头
        setHead(wb, sh);

        HSSFFont f = wb.createFont();
        //字号
        f.setFontHeightInPoints((short) 11);
        //加粗
        f.setBold(true);
        //f.setBoldweight(HSSFFont.BOLDWEIGHT_NORMAL);
        //设置头部字体为宋体
        f.setFontName("宋体");

        HSSFCellStyle style = wb.createCellStyle();
        style.setFont(f);
        //左右居中
        style.setAlignment(HorizontalAlignment.CENTER);
        style.setVerticalAlignment(VerticalAlignment.CENTER);
        style.setFillBackgroundColor(IndexedColors.WHITE.getIndex());
        //填充数据。
        //查询出来所有的员工，不超过24000条的话直接输出了
        List<PortalUser> userList;
        if (StringUtil.isEmpty(deptId)) {
            userList = this.portalUserRepository.findByTenantId(tenantId);
        } else {
            userList = this.portalUserRepository.findPortalUserByDeptIdIn(tenantId, Arrays.asList(deptId));
        }
        if (userList != null) {
            int i = 1;
            for (PortalUser ue : userList) {
                genUserRow(sh, ue, i);
                i++;
            }
        }

        response.reset();
        // 设置response的Header
        response.setContentType("application/octet-stream");
        response.addHeader("Content-Disposition",
                "attachment;filename=" + new String(
                        ("用户列表数据_" + sdf.format(new Timestamp(System.currentTimeMillis())) + ".xls").getBytes("GBK"),
                        "iso8859-1"));
        response.setContentType("application/msexcel;charset=GBK");

        wb.write(response.getOutputStream());

        wb.close();
    }

    /**
     * 生成一条用户数据的行。
     *
     * @param sh    代表一个sheet页。
     * @param ue    用户数据。
     * @param index 行的下标，从1开始.
     */
    private void genUserRow(Sheet sh, PortalUser ue, int index) {
        Row row = sh.createRow(index);
        row.setHeight((short) 400);
        //登录名
        Cell cell = row.createCell(0);
        cell.setCellValue(ue.getLoginId());
        //姓名
        cell = row.createCell(1);
        cell.setCellValue(ue.getName());
        //部门名称（全名）
        cell = row.createCell(2);

        cell.setCellValue(this.getDeptFullNameByDeptId(ue.getDeptId()));

        //职务
        cell = row.createCell(3);
        cell.setCellValue(ue.getTitle());
        //密码
        cell = row.createCell(4);
        cell.setCellValue("");
        //邮箱
        cell = row.createCell(5);
        cell.setCellValue(ue.getEmail());
        //移动电话
        cell = row.createCell(6);
        cell.setCellValue(ue.getMobile1());
        //办公电话
        cell = row.createCell(7);
        cell.setCellValue(ue.getOfficeTel1());
        //数据来源
        cell = row.createCell(8);
        cell.setCellValue(SourceTypeUtil.getSourceType(ue.getSourceType()));
        //备注
        cell = row.createCell(9);
        cell.setCellValue(ue.getRemark());
    }

    /**
     * 根据部门Id查询部门fullName
     *
     * @param deptId
     * @return
     */
    private String getDeptFullNameByDeptId(String deptId) {
        if (StringUtil.isEmpty(deptId)) {
            //如果用户没有部门Id，返回空
            return null;
        }
        String childId = deptId;
        StringBuffer deptFullName = new StringBuffer();
        while (!StringUtil.isEmpty(childId)) {
            PortalDept portalDept = portalDeptRepository.findByDeptId(childId);
            if (portalDept == null) {
                childId = null;
                continue;
            }
            deptFullName=new StringBuffer("/").append(portalDept.getDisplayName()).append(deptFullName);
            childId = portalDept.getParentId();
        }
        return StringUtil.isEmpty(deptFullName.toString()) ? null
                : deptFullName.toString();
    }

}