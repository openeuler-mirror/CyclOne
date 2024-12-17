/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web.auto.excelfilehandler;

import com.idcos.cloud.biz.common.check.CommonParamtersChecker;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.cloud.core.common.util.StringUtil;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryTemplate;
import com.idcos.enterprise.portal.biz.common.utils.Base64Util;
import com.idcos.enterprise.portal.biz.common.utils.PasswordUtil;
import com.idcos.enterprise.portal.biz.common.utils.SourceTypeUtil;
import com.idcos.enterprise.portal.dal.entity.PortalDept;
import com.idcos.enterprise.portal.dal.entity.PortalUser;
import com.idcos.enterprise.portal.dal.enums.ArrayNumEnum;
import com.idcos.enterprise.portal.dal.enums.PortalUserStatusEnum;
import com.idcos.enterprise.portal.dal.repository.PortalDeptRepository;
import com.idcos.enterprise.portal.dal.repository.PortalUserRepository;
import com.idcos.enterprise.portal.vo.PortalUserImportVO;
import org.apache.poi.hssf.usermodel.HSSFWorkbook;
import org.apache.poi.poifs.filesystem.POIFSFileSystem;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Row;
import org.apache.poi.ss.usermodel.Sheet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import javax.servlet.http.HttpServletRequest;
import java.io.IOException;
import java.util.*;

import static com.idcos.enterprise.portal.UamConstant.SLASH;
import static com.idcos.enterprise.portal.UamConstant.STAR;

/**
 * @author Dana
 * @version UserExcelFileHandler.ava, v1 2017/12/10 下午6:08 Dana Exp $$
 */
@Service
public class UserExcelFileHandler {
    private static final Logger log = LoggerFactory.getLogger(UserExcelFileHandler.class);

    @Autowired
    private PortalUserRepository portalUserRepository;

    @Autowired
    private PortalDeptRepository portalDeptRepository;

    @Autowired
    private BusinessQueryTemplate businessQueryTemplate;

    /**
     * 用户导入的excel表头模板
     */
    private static final String[] USER_COLUMNS = {"登录名", "姓名", "具体部门名称", "职务", "密码", "邮箱", "移动电话", "办公电话", "数据来源", "备注"};

    /**
     * 当前上传文件的文件名称
     */
    private String fileName;

    /**
     * 当前上传文件的文件名称
     */
    private final String fileType = ".xls";

    /**
     * step1:处理上传文件请求request。
     *
     * @param request
     * @return
     */
    public CommonResult<?> processUploadFile(final HttpServletRequest request, final String tenantId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<String>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(tenantId);
            }

            @Override
            public String doQuery() {
                //获得上传的文件
                MultipartFile file = UploadFileUtils.getFirstUploadFile(request);
                if (file == null) {
                    log.warn("上传的文件为空。");
                    throw new RuntimeException("上传的文件为空。");
                }
                //取得当前上传文件的文件名称
                fileName = file.getOriginalFilename();
                log.info("Received upload file: " + fileName);
                if (StringUtil.isBlank(fileName)) {
                    log.error(String.format("Upload file is empty. "));
                    throw new RuntimeException("上传的文件为空。");
                }
                //检查文件类型，如果后缀不为.xls则报错
                if (!fileName.toLowerCase().endsWith(fileType)) {
                    log.info(String.format("Upload file is not a valid Excel file: ", fileName));
                    throw new RuntimeException("上传的文件不是Excel文件：" + fileName);
                }
                //记录上传过程起始时的时间，用来计算上传时间
                long pre = System.currentTimeMillis();
                try {
                    //-->step2:解析上传的excel数据文件。
                    parseFile(file, tenantId);
                    //如果解析过程中发现内容有错，则直接返回给客户端处理。
                    //TODO
                } catch (IOException e) {
                    log.info("Upload file read error.", e);
                    throw new RuntimeException(String.format("读取上传文件:%s时出错，原因：%s。", fileName, e.getMessage()));
                } catch (Exception e) {
                    log.info("Process data failed", e);
                    throw new RuntimeException(String.format("处理上传文件数据时出错，原因：%s。", e.getMessage()));
                }

                //记录上传该文件后的时间
                long finalTime = System.currentTimeMillis();
                log.info("Parsed  file: {}, elapsed time ：{}ms.", fileName, finalTime - pre);
                return fileName;
            }
        });
    }

    /**
     * step2:解析上传的excel数据文件。
     */
    public void parseFile(MultipartFile file, String tenantId) throws IOException {
        POIFSFileSystem fs = new POIFSFileSystem(file.getInputStream());
        HSSFWorkbook workbook = new HSSFWorkbook(fs);
        //只读取第一个sheet页。
        String sn = workbook.getSheetName(0);
        try {

            log.info("Parsing and process user of {}...", sn);
            Sheet sheet = workbook.getSheetAt(0);
            //-->step3:读取每个sheet中的数据。
            readSheet(sheet, tenantId);
        } catch (ExcelImportException e) {
            log.error(String.format("Parse sheet:%s failed.", sn), e);
        }
        workbook.close();
    }

    /**
     * step3:读取每个sheet中的数据。
     *
     * @param sheet Sheet实例。
     * @return
     * @throws ExcelImportException 如果插入或更新excel时出错抛出此异常。
     */
    private void readSheet(Sheet sheet, String tenantId) throws ExcelImportException {
        //第一行: 列名。
        int firstRowIndex = sheet.getFirstRowNum();
        Row row = sheet.getRow(firstRowIndex);

        //列名的index和名字之间的对应关系为actualColNames。
        Map<Integer, String> actualColNames = new HashMap<>(9);
        for (int i = row.getFirstCellNum(); i <= row.getLastCellNum(); i++) {
            Cell cc = row.getCell(i);
            String value = ExcelUtil.getValueFromCell(cc);
            if (StringUtil.isBlank(value)) {
                continue;
            }
            actualColNames.put(i, value.trim());
        }
        List<PortalUserImportVO> userList = new ArrayList<>();
        for (int i = firstRowIndex + 1; i <= sheet.getLastRowNum(); i++) {
            row = sheet.getRow(i);
            if (row == null) {
                continue;
            }
            //-->step4:处理每一行数据
            PortalUserImportVO user = parseUserRow(row, actualColNames);
            if (user == null) {
                //此行有问题，忽略。
                continue;
            }
            user.setTenantId(tenantId);
            //此处DataStatus该字段值为空时，才能考虑是新增还是更新，否则该行数据有问题
            if (user.getDataStatus() == null) {
                PortalUser tmpUser = portalUserRepository.findPortalUserById(user.getTenantId(), user.getLoginId());
                if (tmpUser == null) {
                    //1代表新增，0代表更新
                    user.setDataStatus("1");
                    //如果是新增，密码为空，则随机生成一个8位密码
                    if (StringUtil.isBlank(user.getPassword())) {
                        String newPassword = PasswordUtil.getStringRandom(8);
                        user.setPassword(newPassword);
                        //DataStatus字段值改为2，说明导入时没有给密码，用户状态只能给init
                        user.setDataStatus("2");
                    }
                } else {
                    user.setDataStatus("0");
                    //如果是更新，同时不想修改密码，给密码一个状态
                    if (StringUtil.isBlank(user.getPassword())) {
                        user.setPassword("*");
                    }
                }
            }
            userList.add(user);
        }
        //如果登录id相同，则提示给前端，不让保存
        Map<?, List<PortalUserImportVO>> listMap = ListUtil.groupBy(userList, "loginId");
        Set<?> keys = listMap.keySet();
        Iterator<?> iterator = keys.iterator();
        while (iterator.hasNext()) {
            String key = (String) iterator.next();
            List<PortalUserImportVO> portalUserImportVOS = listMap.get(key);
            if (portalUserImportVOS.size() > 1) {
                for (PortalUserImportVO portalUserImportVO : portalUserImportVOS) {
                    portalUserImportVO.setDataStatus("登录名不允许重复");
                }
            }
        }
        //将userList存入缓存中
        UserExcelManager.getInstance().addUserImportList(fileName, userList);
    }

    /**
     * step4:处理每一行数据，处理前先验证数据对应的列是否在预定的列中。
     *
     * @param row
     * @param actualColNames
     * @return
     */
    private PortalUserImportVO parseUserRow(Row row, Map<Integer, String> actualColNames) {
        PortalUserImportVO user = new PortalUserImportVO();
        for (Map.Entry<Integer, String> kv : actualColNames.entrySet()) {
            int col = kv.getKey();
            String colName = kv.getValue();
            Cell cell = row.getCell(col);
            //数据对应的列是否在预定的列中。
            boolean exist = false;
            for (int i = 0; i < USER_COLUMNS.length; i++) {
                if (colName.equals(USER_COLUMNS[i])) {
                    exist = true;
                    break;
                }
            }
            if (!exist) {
                log.warn("行：{},列：{}({})不在系统内置表字段中，将忽略。", "" + row.getRowNum(), "" + col, colName);
                continue;
            }
            String value = ExcelUtil.getValueFromCell(cell);
            try {
                //-->step5:处理每一行各个字段（登录名和姓名不可为空）
                processColumn(value, colName, user);
            } catch (Exception ex) {
                log.error(String.format("处理第%s行时出现系统错误: " + ex.getMessage(), row.getRowNum()), ex);
                break;
            }
        }
        return user;
    }

    /**
     * step5:处理每一行各个字段（登录名和姓名不可为空）
     *
     * @param value
     * @param colName
     * @param user
     * @return
     * @throws ExcelImportException
     */
    private void processColumn(String value, String colName, PortalUserImportVO user) throws ExcelImportException {
        StringBuilder sb = new StringBuilder();
        if (colName.equals(USER_COLUMNS[0])) {
            //登录名
            if (StringUtil.isBlank(value)) {
                sb.append("登录名为空；");
                user.setDataStatus(sb.toString());
            } else {
                user.setLoginId(value.trim());
            }
        } else if (colName.equals(USER_COLUMNS[1])) {
            //姓名
            if (StringUtil.isBlank(value)) {
                sb.append("姓名为空；");
                user.setDataStatus(sb.toString());
            } else {
                user.setName(value.trim());
            }
        } else if (colName.equals(USER_COLUMNS[ArrayNumEnum.TWO.getCode()])) {
            //具体部门名称
            if (StringUtil.isBlank(value)) {
                user.setDeptFullName(null);
            } else {
                if (!value.startsWith(SLASH)) {
                    user.setDataStatus("部门名称必须以'/'开头");
                }
                String[] deptNameArray = value.split("/");
                for (int i = 1; i < deptNameArray.length; i++) {
                    if (StringUtil.isBlank(deptNameArray[i])) {
                        user.setDataStatus("部门名称格式错误");
                        break;
                    }
                }
                user.setDeptFullName(value.trim());
            }
        } else if (colName.equals(USER_COLUMNS[ArrayNumEnum.THREE.getCode()])) {
            //职务
            if (StringUtil.isBlank(value)) {
                user.setTitle(null);
            } else {
                user.setTitle(value.trim());
            }
        } else if (colName.equals(USER_COLUMNS[ArrayNumEnum.FOUR.getCode()])) {
            //密码
            if (StringUtil.isBlank(value)) {
                user.setPassword(null);
            } else {
                try {
                    PasswordUtil.checkPassword(value);
                } catch (Exception e) {
                    user.setDataStatus(e.getMessage());
                }
                user.setPassword(value.trim());
            }
        } else if (colName.equals(USER_COLUMNS[ArrayNumEnum.FIVE.getCode()])) {
            //邮箱
            if (StringUtil.isBlank(value)) {
                user.setEmail(null);
            } else {
                user.setEmail(value.trim());
            }
        } else if (colName.equals(USER_COLUMNS[ArrayNumEnum.SIX.getCode()])) {
            //移动电话
            if (StringUtil.isBlank(value)) {
                user.setMobile1(null);
            } else {
                user.setMobile1(value.trim());
            }
        } else if (colName.equals(USER_COLUMNS[ArrayNumEnum.SEVEN.getCode()])) {
            //办公电话
            if (StringUtil.isBlank(value)) {
                user.setOfficeTel1(null);
            } else {
                user.setOfficeTel1(value.trim());
            }
        } else if (colName.equals(USER_COLUMNS[ArrayNumEnum.EIGHT.getCode()])) {
            //数据来源
            user.setSourceType(SourceTypeUtil.setSourceType(value.trim()));
        } else if (colName.equals(USER_COLUMNS[ArrayNumEnum.NINE.getCode()])) {
            //备注
            if (StringUtil.isBlank(value)) {
                user.setRemark(null);
            } else {
                user.setRemark(value.trim());
            }
        } else {
        }

    }

    /**
     * 根据excel名，返回上传预览。
     *
     * @param fileName
     * @return
     */
    public CommonResult<?> getUsersByExcelFileName(final String fileName) {
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalUserImportVO>>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(fileName);
            }

            @Override
            public List<PortalUserImportVO> doQuery() {
                //从缓存中取userList
                List<PortalUserImportVO> userImportList = UserExcelManager.getInstance().getUserImportList(fileName);
                return userImportList;
            }
        });
    }

    /**
     * 导入预览中的合法数据，并删除用户列表缓存。
     *
     * @param fileName
     * @return
     */
    public CommonResult<?> saveUsersByExcelFileName(final String fileName) {
        return businessQueryTemplate.process(new BusinessQueryCallback<String>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(fileName);
            }

            @Override
            public String doQuery() {
                long beginTime = System.currentTimeMillis();
                log.info("开始导入数据: " + beginTime);
                //从缓存中取出userImportList，将dataStatus值为空的数据存入数据库
                List<PortalUserImportVO> userImportList = UserExcelManager.getInstance().getUserImportList(fileName);
                //需要插入的数据
                List<PortalUser> insertList = new ArrayList<>();
                //需要更新的数据
                List<PortalUser> updateList = new ArrayList<>();
                for (PortalUserImportVO userImportVO : userImportList) {
                    if (userImportVO.getDataStatus() != "0" && userImportVO.getDataStatus() != "1"
                            && userImportVO.getDataStatus() != "2") {
                        continue;
                    }
                    PortalUser tmpUser = portalUserRepository.findPortalUserById(userImportVO.getTenantId(),
                            userImportVO.getLoginId());
                    if (tmpUser == null) {
                        //插入
                        PortalUser insertUser = getInsertingUser(userImportVO);
                        insertList.add(insertUser);
                    } else {
                        //更新
                        PortalUser updateUser = getUpdatingUser(userImportVO, tmpUser);
                        if (updateUser == null) {
                            continue;
                        }
                        updateList.add(updateUser);
                    }
                }

                List<PortalUser> insertUsers = portalUserRepository.save(insertList);
                portalUserRepository.save(updateList);

                producePassword(insertUsers);

                //从缓存中删除userList
                UserExcelManager.getInstance().removeUserImportList(fileName);
                long endTime = System.currentTimeMillis();
                long duringTime = endTime - beginTime;
                log.info("导入数据结束: " + endTime + ", 共耗时：" + duringTime);
                return "success";
            }
        });
    }

    private PortalUser getUpdatingUser(PortalUserImportVO userImportVO, PortalUser tmpUser) {
        PortalUser updateUser = new PortalUser();
        updateUser.setTenantId(userImportVO.getTenantId());
        updateUser.setId(tmpUser.getId());
        updateUser.setLoginId(userImportVO.getLoginId());
        updateUser.setName(userImportVO.getName());
        updateUser.setTitle(userImportVO.getTitle());
        updateUser.setEmail(userImportVO.getEmail());
        updateUser.setMobile1(userImportVO.getMobile1());
        updateUser.setOfficeTel1(userImportVO.getOfficeTel1());
        updateUser.setSourceType(userImportVO.getSourceType());
        updateUser.setRemark(userImportVO.getRemark());
        updateUser.setIsActive("Y");
        updateUser.setCreateUser(tmpUser.getCreateUser());
        updateUser.setCreateTime(tmpUser.getCreateTime());
        updateUser.setLastModifiedTime(new Date());
        updateUser.setStatus(tmpUser.getStatus());

        try {
            //如果上传文件中没有密码，则不对密码进行修改。
            if (STAR.equals(userImportVO.getPassword())) {
                //将原来的密码放到此对象中。
                updateUser.setPassword(tmpUser.getPassword());
                updateUser.setSalt(tmpUser.getSalt());
            } else { //否则需要修改密码
                //salt随机产生，密码是用户的ID.
                byte[] salt = PasswordUtil.getSalt();
                String encriptPW = PasswordUtil.encrypt(userImportVO.getPassword(), updateUser.getId(),
                        salt);
                updateUser.setSalt(Base64Util.encode(salt));
                updateUser.setPassword(encriptPW);
            }
        } catch (Exception ex) {
            return null;
        }
        String deptId = parseDeptFullName(userImportVO.getDeptFullName(), userImportVO.getTenantId(), userImportVO.getSourceType());
        updateUser.setDeptId(deptId);
        return updateUser;
    }

    private PortalUser getInsertingUser(PortalUserImportVO userImportVO) {
        PortalUser insertUser = new PortalUser();
        insertUser.setTenantId(userImportVO.getTenantId());
        insertUser.setLoginId(userImportVO.getLoginId());
        insertUser.setName(userImportVO.getName());
        insertUser.setTitle(userImportVO.getTitle());
        insertUser.setEmail(userImportVO.getEmail());
        insertUser.setMobile1(userImportVO.getMobile1());
        insertUser.setOfficeTel1(userImportVO.getOfficeTel1());
        insertUser.setSourceType(userImportVO.getSourceType());
        insertUser.setRemark(userImportVO.getRemark());
        insertUser.setIsActive("Y");
        insertUser.setCreateUser("import");
        insertUser.setCreateTime(new Date());
        insertUser.setLastModifiedTime(new Date());
        if (userImportVO.getDataStatus() == "1") {
            insertUser.setStatus(PortalUserStatusEnum.ENABLED.getCode());
        } else if (userImportVO.getDataStatus() == "2") {
            insertUser.setStatus(PortalUserStatusEnum.INIT.getCode());
        }
        insertUser.setPassword(userImportVO.getPassword());
        String deptId = parseDeptFullName(userImportVO.getDeptFullName(), userImportVO.getTenantId(), userImportVO.getSourceType());
        insertUser.setDeptId(deptId);
        return insertUser;
    }

    private void producePassword(List<PortalUser> users) {
        List<PortalUser> insertUsers = new ArrayList<>();
        for (PortalUser user : users) {
            try {
                byte[] salt = PasswordUtil.getSalt();
                String encriptPW = PasswordUtil.encrypt(user.getPassword(), user.getId(), salt);
                user.setSalt(Base64Util.encode(salt));
                user.setPassword(encriptPW);
                insertUsers.add(user);
            } catch (Exception ex) {
                continue;
            }
        }
        portalUserRepository.save(insertUsers);
    }


    /**
     * 处理部门全名，返回部门id
     *
     * @param deptFullName
     * @param tenantId
     */
    private String parseDeptFullName(String deptFullName, String tenantId, String sourceType) {
        if (StringUtil.isBlank(deptFullName)) {
            return null;
        }
        //以符号/给部门名称分段，爷/父/子
        String[] deptNameArray = deptFullName.substring(1).split("/");
        PortalDept parentDept = portalDeptRepository.findByDeptNameAndTenantId(deptNameArray[0], tenantId);
        int arrayNo = 0;
        return produceDeptId(parentDept, tenantId, deptNameArray, arrayNo, "", sourceType);
    }

    private String produceDeptId(PortalDept parentDept, String tenantId, String[] deptNameArray, int arrayNo,
                                 String parentId, String sourceType) {
        if (parentDept == null) {
            return createRecurseDept(tenantId, deptNameArray, parentId, arrayNo, sourceType);
        }
        if (arrayNo < deptNameArray.length - 1) {
            PortalDept dept = portalDeptRepository.findByDeptNameAndParentIdAndTenantId(deptNameArray[arrayNo + 1],
                    tenantId, parentDept.getId());
            arrayNo++;
            return produceDeptId(dept, tenantId, deptNameArray, arrayNo, parentDept.getId(), sourceType);
        }
        return parentDept.getId();
    }

    private String createRecurseDept(String tenantId, String[] deptNameArray, String parentId, int arrayNo, String sourceType) {
        //级联创建爷/父/子
        if (arrayNo < deptNameArray.length) {
            PortalDept parentDept = new PortalDept();
            parentDept.setCode(deptNameArray[arrayNo]);
            parentDept.setDisplayName(deptNameArray[arrayNo]);
            parentDept.setTenantId(tenantId);
            parentDept.setStatus("1");
            parentDept.setGmtCreate(new Date());
            parentDept.setGmtModified(new Date());
            parentDept.setParentId(parentId);
            parentDept.setSourceType(sourceType);
            PortalDept dbparentDept = portalDeptRepository.save(parentDept);
            arrayNo++;
            return createRecurseDept(tenantId, deptNameArray, dbparentDept.getId(), arrayNo, sourceType);
        }
        return parentId;
    }

    /**
     * 取消导入，直接删除用户列表缓存
     *
     * @param fileName
     * @return
     */
    public CommonResult<?> removeUsersByExcelFileName(final String fileName) {
        return businessQueryTemplate.process(new BusinessQueryCallback<String>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(fileName);
            }

            @Override
            public String doQuery() {
                //从缓存中删除userList
                UserExcelManager.getInstance().removeUserImportList(fileName);
                return "success";
            }
        });
    }

}