<template>
  <router-view />
</template>

<script setup>
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

#app {
  height: 100vh;
  overflow: hidden;
}

.el-card {
  border-radius: 12px;
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
}

.el-card:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.el-card .el-card__header {
  border-bottom: 1px solid #f0f0f0;
  padding: 18px 20px;
  background: #fafbfc;
}

.el-card .el-card__body {
  padding: 20px;
}

.el-button--primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  padding: 10px 20px;
  font-weight: 500;
  border-radius: 6px;
  transition: all 0.3s ease;
}

.el-button--primary:hover {
  background: linear-gradient(135deg, #5568d3 0%, #6a3f9c 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.el-button--primary:active {
  transform: translateY(0);
}

.el-table {
  border-radius: 8px;
  overflow: hidden;
  border: none;
}

.el-table th {
  background: #764ba2 !important;
  color: white !important;
  font-weight: 600;
  padding: 14px 0;
}

.el-table td {
  padding: 12px 0;
  background: white !important;
}

.el-table .el-table__row:hover > td {
  background-color: #f8f9fc !important;
}

.el-table .el-table__fixed td {
  background: white !important;
}

.el-table .el-table__fixed-right td {
  background: white !important;
}

.el-pagination {
  margin-top: 20px;
}

.el-pagination .el-pager li {
  border-radius: 6px;
  font-weight: 500;
}

.el-pagination .el-pager li.is-active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.el-input__wrapper {
  border-radius: 6px;
  transition: all 0.3s ease;
}

.el-input__wrapper:hover {
  box-shadow: 0 0 0 1px #667eea inset;
}

.el-input__wrapper.is-focus {
  box-shadow: 0 0 0 1px #667eea inset;
}

.el-dialog {
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
}

.el-dialog__header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px 24px;
  border-radius: 12px 12px 0 0;
}

.el-dialog__title {
  color: white;
  font-weight: 600;
  font-size: 18px;
}

.el-dialog__headerbtn .el-dialog__close {
  color: white;
  font-size: 20px;
}

.el-dialog__headerbtn:hover .el-dialog__close {
  color: rgba(255, 255, 255, 0.8);
}

.el-dialog__body {
  padding: 24px;
}

.el-form-item__label {
  font-weight: 500;
  color: #333;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .el-card .el-card__header {
    padding: 14px 16px;
  }

  .el-card .el-card__body {
    padding: 16px;
  }

  .el-dialog {
    width: 90% !important;
    margin: 0 auto;
  }

  .el-dialog__header {
    padding: 16px 20px;
  }

  .el-dialog__body {
    padding: 16px 20px;
  }

  .el-dialog__title {
    font-size: 16px;
  }

  /* 表格滚动 */
  .el-table {
    font-size: 12px;
  }

  .el-table th,
  .el-table td {
    padding: 8px 4px;
  }

  .el-table .cell {
    padding: 0 4px;
  }

  /* 分页器 */
  .el-pagination {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
  }

  .el-pagination .el-pager li {
    min-width: 28px;
    height: 28px;
    font-size: 12px;
  }

  .el-pagination button {
    min-width: 28px;
    height: 28px;
  }

  /* 表单 */
  .el-form-item {
    margin-bottom: 16px;
  }

  .el-form-item__label {
    font-size: 13px;
  }

  /* 按钮 */
  .el-button {
    padding: 8px 12px;
    font-size: 13px;
  }

  /* 输入框 */
  .el-input__inner {
    font-size: 14px;
  }

  /* 卡片响应式 */
  .el-col {
    margin-bottom: 12px;
  }
}

@media (max-width: 480px) {
  .el-card .el-card__header {
    padding: 12px;
  }

  .el-card .el-card__body {
    padding: 12px;
  }

  .el-dialog {
    width: 95% !important;
  }

  .el-dialog__body {
    padding: 12px;
  }

  .el-table {
    font-size: 11px;
  }

  /* 超小屏幕隐藏部分列 */
  .el-table .small-hide {
    display: none;
  }
}
</style>
