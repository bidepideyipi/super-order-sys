<template>
  <el-container class="app-container">
    <el-header class="app-header">
      <div class="header-left">
        <el-button
          class="menu-toggle"
          :icon="Menu"
          @click="drawerVisible = true"
          circle
        />
        <h1>超级订单管理系统</h1>
      </div>
      <div class="header-right">
        <el-button type="info" plain @click="handleLogout">退出</el-button>
      </div>
    </el-header>

    <el-container>
      <!-- PC端侧边栏 -->
      <el-aside width="200px" class="app-aside">
        <el-menu
          :default-active="activeMenu"
          @select="handleMenuSelect"
          class="sidebar-menu"
        >
          <el-menu-item index="/">
            <el-icon><HomeFilled /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="/sku">
            <el-icon><Goods /></el-icon>
            <span>SKU 管理</span>
          </el-menu-item>
          <el-menu-item index="/orders">
            <el-icon><Document /></el-icon>
            <span>订单管理</span>
          </el-menu-item>
          <el-menu-item index="/purchases">
            <el-icon><ShoppingCart /></el-icon>
            <span>采购管理</span>
          </el-menu-item>
          <el-menu-item index="/settlement">
            <el-icon><Coin /></el-icon>
            <span>结算管理</span>
          </el-menu-item>
          <el-menu-item index="/financial">
            <el-icon><Wallet /></el-icon>
            <span>财务管理</span>
          </el-menu-item>
          <el-menu-item index="/customers">
            <el-icon><User /></el-icon>
            <span>客户管理</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>

    <!-- 移动端抽屉菜单 -->
    <el-drawer
      v-model="drawerVisible"
      :size="280"
      direction="ltr"
      :with-header="false"
    >
      <el-menu
        :default-active="activeMenu"
        @select="handleMobileMenuSelect"
        class="mobile-menu"
      >
        <el-menu-item index="/">
          <el-icon><HomeFilled /></el-icon>
          <span>首页</span>
        </el-menu-item>
        <el-menu-item index="/sku">
          <el-icon><Goods /></el-icon>
          <span>SKU 管理</span>
        </el-menu-item>
        <el-menu-item index="/orders">
          <el-icon><Document /></el-icon>
          <span>订单管理</span>
        </el-menu-item>
        <el-menu-item index="/purchases">
          <el-icon><ShoppingCart /></el-icon>
          <span>采购管理</span>
        </el-menu-item>
        <el-menu-item index="/settlement">
          <el-icon><Coin /></el-icon>
          <span>结算管理</span>
        </el-menu-item>
        <el-menu-item index="/financial">
          <el-icon><Wallet /></el-icon>
          <span>财务管理</span>
        </el-menu-item>
        <el-menu-item index="/customers">
          <el-icon><User /></el-icon>
          <span>客户管理</span>
        </el-menu-item>
      </el-menu>
    </el-drawer>
  </el-container>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { HomeFilled, Goods, Document, User, ShoppingCart, Wallet, Coin, Menu } from '@element-plus/icons-vue';
import { useAuthStore } from '@/stores/auth';
import { ElMessage, ElMessageBox } from 'element-plus';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const activeMenu = ref(route.path);
const drawerVisible = ref(false);

watch(() => route.path, (newPath) => {
  activeMenu.value = newPath;
}, { immediate: true });

const handleMenuSelect = (index) => {
  activeMenu.value = index;
  router.push(index).catch(err => {
    console.error('Router push error:', err);
  });
};

const handleMobileMenuSelect = (index) => {
  activeMenu.value = index;
  drawerVisible.value = false;
  router.push(index).catch(err => {
    console.error('Router push error:', err);
  });
};

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    authStore.logout();
    ElMessage.success('已退出登录');
    router.push('/login');
  } catch {
    // 用户取消
  }
};
</script>

<style scoped>
.app-container {
  height: 100vh;
}

.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 0 24px;
  height: 60px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-left h1 {
  font-size: 22px;
  font-weight: 600;
  letter-spacing: 0.5px;
  background: linear-gradient(90deg, rgba(255,255,255,0.95), rgba(255,255,255,0.85));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  white-space: nowrap;
}

.menu-toggle {
  display: none;
}

.app-aside {
  background: #ffffff;
  border-right: 1px solid #e8eaec;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.04);
}

.sidebar-menu {
  border-right: none;
  padding: 8px 0;
}

.sidebar-menu .el-menu-item,
.mobile-menu .el-menu-item {
  margin: 4px 12px;
  border-radius: 8px;
  height: 44px;
  line-height: 44px;
  color: #606266;
  transition: all 0.3s ease;
}

.sidebar-menu .el-menu-item:hover,
.mobile-menu .el-menu-item:hover {
  background-color: #f5f7fa;
  color: #667eea;
  transform: translateX(4px);
}

.sidebar-menu .el-menu-item.is-active,
.mobile-menu .el-menu-item.is-active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 2px 6px rgba(102, 126, 234, 0.3);
}

.sidebar-menu .el-menu-item .el-icon,
.mobile-menu .el-menu-item .el-icon {
  margin-right: 8px;
  font-size: 18px;
}

.sidebar-menu .el-menu-item span,
.mobile-menu .el-menu-item span {
  font-weight: 500;
  font-size: 14px;
}

.app-main {
  padding: 24px;
  background: linear-gradient(180deg, #f8f9fc 0%, #f0f2f5 100%);
  overflow-y: auto;
}

.mobile-menu {
  border-right: none;
  padding: 16px 0;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .menu-toggle {
    display: flex;
  }

  .header-left h1 {
    font-size: 16px;
  }

  .app-aside {
    display: none;
  }

  .app-main {
    padding: 16px;
  }

  .header-right .el-button {
    padding: 8px 12px;
    font-size: 12px;
  }
}

@media (max-width: 480px) {
  .app-header {
    padding: 0 12px;
  }

  .header-left h1 {
    font-size: 14px;
  }

  .app-main {
    padding: 12px;
  }
}
</style>
