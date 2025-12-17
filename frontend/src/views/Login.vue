<template>
  <div class="login-container">
    <el-card class="login-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>电力监控系统工作报备工具</span>
        </div>
      </template>
      <el-form
        :model="loginForm"
        status-icon
        :rules="rules"
        ref="loginFormRef"
        label-width="80px"
        class="login-form"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
          ></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            @click="handleLogin"
            :loading="loading"
            style="width: 100%"
            >登录</el-button
          >
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useStore } from "vuex";
import { ElMessage } from "element-plus";
import { login } from "@/api/user";

export default {
  name: "Login",
  setup() {
    const router = useRouter();
    const store = useStore();
    const loginFormRef = ref();
    const loading = ref(false);

    const loginForm = ref({
      username: "",
      password: "",
    });

    const rules = {
      username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
      password: [{ required: true, message: "请输入密码", trigger: "blur" }],
    };

    const handleLogin = async () => {
      if (!loginFormRef.value) return;

      await loginFormRef.value.validate(async (valid) => {
        if (!valid) return;

        loading.value = true;

        try {
          console.log("发送登录请求:", {
            username: loginForm.value.username,
            password: loginForm.value.password,
          });

          const response = await login({
            username: loginForm.value.username,
            password: loginForm.value.password,
          });

          console.log("登录响应:", response.data);

          // 提取后端返回的数据
          const { token, user } = response.data;

          if (!token) {
            throw new Error("服务器未返回token");
          }

          // 准备用户信息对象（将后端字段映射到前端字段）
          const userInfo = {
            username: user.username,
            role: user.role,
            userId: user.user_id, // 后端是 user_id，前端是 userId
            realName: user.real_name, // 后端是 real_name，前端是 realName
            unit: user.unit,
            phone: user.phone,
            position: user.position,
          };

          // 调用Vuex action，传入真实token和用户信息
          await store.dispatch("login", {
            token,
            userInfo,
          });

          console.log("登录成功，Vuex状态:", {
            token: store.state.token,
            userInfo: store.state.userInfo,
            isLoggedIn: store.getters.isLoggedIn,
            userRole: store.getters.userRole,
          });

          ElMessage.success("登录成功");

          // 稍微延迟一下确保状态更新
          await new Promise((resolve) => setTimeout(resolve, 50));

          // 根据角色跳转
          const role = store.getters.userRole || user.role;
          console.log("跳转前的角色:", role);

          if (role === "admin") {
            router.push("/admin/approve-report");
          } else {
            router.push("/user/create-report");
          }
        } catch (error) {
          console.error("登录失败:", error);
          console.error("错误详情:", error.response?.data);
          ElMessage.error(
            error.response?.data?.error || error.message || "登录失败"
          );
        } finally {
          loading.value = false;
        }
      });
    };

    return {
      loginForm,
      rules,
      loginFormRef,
      handleLogin,
      loading,
    };
  },
};
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f0f2f5;
}
.login-card {
  width: 400px;
}
.card-header {
  text-align: center;
  font-size: 18px;
  font-weight: bold;
}
.login-form {
  margin-top: 20px;
}
</style>