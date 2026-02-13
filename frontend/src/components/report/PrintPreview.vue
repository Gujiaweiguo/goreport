<template>
  <el-dialog
    v-model="visible"
    title="打印预览"
    width="600px"
    :before-close="handleClose"
  >
    <el-form :model="settings" label-width="120px">
      <el-form-item label="纸张大小">
        <el-select v-model="settings.paperSize">
          <el-option label="A4 (210mm × 297mm)" value="A4" />
          <el-option label="A3 (297mm × 420mm)" value="A3" />
          <el-option label="Letter (216mm × 279mm)" value="Letter" />
        </el-select>
      </el-form-item>

      <el-form-item label="页面方向">
        <el-radio-group v-model="settings.orientation">
          <el-radio label="portrait">纵向</el-radio>
          <el-radio label="landscape">横向</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="页边距 (mm)">
        <el-row :gutter="10">
          <el-col :span="12">
            <el-form-item label="上边距" label-width="70px">
              <el-input-number v-model="settings.marginTop" :min="0" :max="50" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="下边距" label-width="70px">
              <el-input-number v-model="settings.marginBottom" :min="0" :max="50" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="10">
          <el-col :span="12">
            <el-form-item label="左边距" label-width="70px">
              <el-input-number v-model="settings.marginLeft" :min="0" :max="50" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="右边距" label-width="70px">
              <el-input-number v-model="settings.marginRight" :min="0" :max="50" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form-item>

      <el-form-item>
        <div class="print-preview-area">
          <div 
            class="print-preview-sheet" 
            :style="printPreviewStyle"
          >
            <div class="print-preview-content" v-html="reportHTML"></div>
          </div>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">关闭</el-button>
      <el-button type="primary" @click="handlePrint">打印</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
  import { computed, reactive } from 'vue'

  interface PrintSettings {
    paperSize: string
    orientation: 'portrait' | 'landscape'
    marginTop: number
    marginBottom: number
    marginLeft: number
    marginRight: number
  }

  const props = defineProps<{
    visible: boolean
    reportHTML: string
  }>()

  const emit = defineEmits<{
    'update:visible': [value: boolean]
  }>()

  const settings = reactive<PrintSettings>({
    paperSize: 'A4',
    orientation: 'portrait',
    marginTop: 10,
    marginBottom: 10,
    marginLeft: 15,
    marginRight: 15
  })

  const paperSizeDimensions = {
    A4: { width: 210, height: 297 },
    A3: { width: 297, height: 420 },
    Letter: { width: 216, height: 279 }
  }

  const printPreviewStyle = computed(() => {
    const dim = paperSizeDimensions[settings.paperSize as keyof typeof paperSizeDimensions]
    const isLandscape = settings.orientation === 'landscape'

    return {
      width: `${dim.height}mm`,
      height: `${dim.width}mm`,
      padding: `${settings.marginTop}mm ${settings.marginRight}mm ${settings.marginBottom}mm ${settings.marginLeft}mm`,
      background: '#f5f5f5',
      border: '1px solid #ddd',
      transform: isLandscape ? 'rotate(-90deg)' : 'none',
      transformOrigin: 'center center',
      overflow: 'hidden'
    }
  })

  function handleClose() {
    emit('update:visible', false)
  }

  function handlePrint() {
    const style = document.createElement('style')
    const s = settings

    let css = '@page { '
    css += `size: ${s.paperSize} ${s.orientation}; `
    css += `margin: ${s.marginTop}mm ${s.marginRight}mm ${s.marginBottom}mm ${s.marginLeft}mm; `
    css += '} '

    css += '@media print { '
    css += 'body * { visibility: hidden; } '
    css += '.print-area, .print-area * { visibility: visible; } '
    css += 'body { margin: 0; } '
    css += '}'

    style.textContent = css
    document.head.appendChild(style)

    setTimeout(() => {
      window.print()
      style.remove()
    }, 100)

    emit('update:visible', false)
  }
</script>

<style scoped>
.print-preview-area {
  width: 100%;
  height: 400px;
  background: #f0f2f5;
  border: 1px dashed #d1d5db;
  border-radius: 8px;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
}

.print-preview-sheet {
  background: white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.print-preview-content {
  overflow: hidden;
  height: 100%;
}

.print-preview-content :deep(table) {
  border-collapse: collapse;
  font-size: 10px;
}

.print-preview-content :deep(td) {
  border: 1px solid #ddd;
  padding: 2px 4px;
}
</style>
