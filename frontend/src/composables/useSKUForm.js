import { ref } from 'vue';
import { ElMessage } from 'element-plus';

export function useSKUForm() {
  const dialogVisible = ref(false);
  const dialogMode = ref('add');
  const form = ref({
    sku_code: '',
    name: '',
    category_id: '',
    unit: '个',
    box_spec: '',
    box_quantity: 1,
    cost_price: 0,
    sale_price: 0,
    spec: '',
    image_path: '',
    image_file: null
  });

  const resetForm = () => {
    form.value = {
      sku_code: '',
      name: '',
      category_id: '',
      unit: '个',
      box_spec: '',
      box_quantity: 1,
      cost_price: 0,
      sale_price: 0,
      spec: '',
      image_path: '',
      image_file: null
    };
  };

  const openAddDialog = () => {
    dialogMode.value = 'add';
    resetForm();
    dialogVisible.value = true;
  };

  const openEditDialog = (row) => {
    dialogMode.value = 'edit';
    form.value = { ...row };
    dialogVisible.value = true;
  };

  const handleSave = async () => {
    try {
      const imageBase64 = form.value.image_file || null;
      const skuData = {
        name: form.value.name,
        category_id: form.value.category_id,
        unit: form.value.unit,
        spec: form.value.spec,
        box_spec: form.value.box_spec,
        box_quantity: form.value.box_quantity,
        cost_price: form.value.cost_price,
        sale_price: form.value.sale_price
      };

      // 编辑时需要传递sku_code
      if (dialogMode.value === 'edit') {
        skuData.sku_code = form.value.sku_code;
      }

      if (dialogMode.value === 'add') {
        await window.tauriAPI.sku.create(skuData, imageBase64);
        ElMessage.success('新增成功');
      } else {
        await window.tauriAPI.sku.update(String(form.value.id), skuData, imageBase64);
        ElMessage.success('更新成功');
      }
      dialogVisible.value = false;
    } catch (error) {
      ElMessage.error('保存失败');
      console.error(error);
      throw error;
    }
  };

  return {
    dialogVisible,
    dialogMode,
    form,
    openAddDialog,
    openEditDialog,
    handleSave
  };
}
