<template>
  <el-dialog
    v-model="visible"
    title="导出配置"
    width="600px"
    :before-close="handleClose"
  >
    <el-form :model="options" label-width="140px">
      <el-form-item label="导出格式">
        <el-select v-model="format">
          <el-option label="Excel">
            <el-icon class="option-icon"><document /></el-icon>
          </el-option>
          <el-option label="PDF">
            <el-icon class="option-icon"><document-copy /></el-icon>
          </el-option>
          <el-option label="Word">
            <el-icon class="option-icon"><files /></el-icon>
          </el-option>
          <el-option label="图片">
            <el-icon class="option-icon"><picture /></el-icon>
          </el-option>
        </el-select>
      </el-form-item>

      <el-divider />

      <el-form-item v-if="format === 'excel'" label="包含表头">
        <el-switch v-model="excelOptions.includeHeaders" active-text="是" inactive-text="否" />
      </el-form-item>

      <el-form-item v-if="format === 'excel'" label="包含样式">
        <el-switch v-model="excelOptions.includeStyles" active-text="是" inactive-text="否" />
      </el-form-item>

      <el-form-item v-if="format === 'excel'" label="工作表名称">
        <el-input v-model="excelOptions.sheetName" placeholder="Sheet1" />
      </el-form-item>

      <el-form-item v-if="format === 'excel'" label="编码格式">
        <el-select v-model="excelOptions.encoding">
          <el-option label="UTF-8" value="utf-8" />
          <el-option label="GBK" value="gbk" />
          <el-option label="GB2312" value="gb2312" />
        </el-select>
      </el-form-item>

      <el-form-item v-if="format === 'pdf'" label="纸张大小">
        <el-select v-model="pdfOptions.paperSize">
          <el-option label="A4 (210mm × 297mm)" value="A4" />
          <el-option label="A3 (297mm × 420mm)" value="A3" />
          <el-option label="Letter (216mm × 279mm)" value="Letter" />
        </el-select>
      </el-form-item>

      <el-form-item v-if="format === 'pdf'" label="页面方向">
        <el-radio-group v-model="pdfOptions.orientation">
          <el-radio label="portrait">纵向</el-radio>
          <el-radio label="landscape">横向</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item v-if="format === 'pdf'" label="页边距 (mm)">
        <el-row :gutter="10">
          <el-col :span="12">
            <el-form-item label="上边距" label-width="70px">
              <el-input-number v-model="pdfOptions.marginTop" :min="0" :max="50" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="下边距" label-width="70px">
              <el-input-number v-model="pdfOptions.marginBottom" :min="0" :max="50" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="10">
          <el-col :span="12">
            <el-form-item label="左边距" label-width="70px">
              <el-input-number v-model="pdfOptions.marginLeft" :min="0" :max="50" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="右边距" label-width="70px">
              <el-input-number v-model="pdfOptions.marginRight" :min="0" :max="50" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form-item>

      <el-form-item v-if="format === 'pdf'" label="页码">
        <el-switch v-model="pdfOptions.includePageNumbers" active-text="包含" inactive-text="不包含" />
      </el-form-item>

      <el-form-item v-if="format === 'pdf'" label="导出质量">
        <el-radio-group v-model="pdfOptions.quality">
          <el-radio label="high">高质量</el-radio>
          <el-radio label="medium">标准</el-radio>
          <el-radio label="low">快速</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item v-if="format === 'word'" label="包含表头">
        <el-switch v-model="wordOptions.includeHeaders" active-text="是" inactive-text="否" />
      </el-form-item>

      <el-form-item v-if="format === 'word'" label="包含样式">
        <el-switch v-model="wordOptions.includeStyles" active-text="是" inactive-text="否" />
      </el-form-item>

      <el-form-item v-if="format === 'word'" label="文档模板">
        <el-select v-model="wordOptions.template">
          <el-option label="标准模板" value="standard" />
          <el-option label="简洁模板" value="simple" />
          <el-option label="专业模板" value="professional" />
        </el-select>
      </el-form-item>

      <el-form-item v-if="format === 'image'" label="图片格式">
        <el-radio-group v-model="imageOptions.format">
          <el-radio label="png">PNG</el-radio>
          <el-radio label="jpg">JPG</el-radio>
          <el-radio label="svg">SVG</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item v-if="format === 'image'" label="分辨率 (DPI)">
        <el-input-number v-model="imageOptions.resolution" :min="72" :max="300" />
      </el-form-item>

      <el-form-item v-if="format === 'image'" label="图片质量">
        <el-slider v-model="imageOptions.quality" :min="1" :max="100" />
      </el-form-item>

      <el-form-item v-if="format === 'image'" label="背景颜色">
        <el-radio-group v-model="imageOptions.background">
          <el-radio label="transparent">透明</el-radio>
          <el-radio label="white">白色</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleExport">
        导出
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
  import { computed, reactive, ref, watch } from 'vue'
  import { Document, DocumentCopy, Files, Picture } from '@element-plus/icons-vue'

  interface ExcelOptions {
    includeHeaders: boolean
    includeStyles: boolean
    sheetName: string
    encoding: string
  }

  interface PDFOptions {
    paperSize: string
    orientation: string
    marginTop: number
    marginBottom: number
    marginLeft: number
    marginRight: number
    includePageNumbers: boolean
    quality: string
  }

  interface WordOptions {
    includeHeaders: boolean
    includeStyles: boolean
    template: string
  }

  interface ImageOptions {
    format: string
    resolution: number
    quality: number
    background: string
  }

  const props = defineProps<{
    visible: boolean
  }>()

  const emit = defineEmits<{
    'update:visible': [value: boolean]
    'export': [format: string, options: any]
  }>()

  const format = ref('excel')
  const exporting = ref(false)

  const excelOptions = reactive<ExcelOptions>({
    includeHeaders: true,
    includeStyles: true,
    sheetName: 'Sheet1',
    encoding: 'utf-8'
  })

  const pdfOptions = reactive<PDFOptions>({
    paperSize: 'A4',
    orientation: 'portrait',
    marginTop: 10,
    marginBottom: 10,
    marginLeft: 15,
    marginRight: 15,
    includePageNumbers: false,
    quality: 'medium'
  })

  const wordOptions = reactive<WordOptions>({
    includeHeaders: true,
    includeStyles: true,
    template: 'standard'
  })

  const imageOptions = reactive<ImageOptions>({
    format: 'png',
    resolution: 96,
    quality: 90,
    background: 'transparent'
  })

  const options = computed(() => {
    switch (format.value) {
      case 'excel':
        return excelOptions
      case 'pdf':
        return pdfOptions
      case 'word':
        return wordOptions
      case 'image':
        return imageOptions
      default:
        return excelOptions
    }
  })

  watch(format, (newVal) => {
    if (newVal === 'excel') {
      Object.assign(excelOptions, {
        includeHeaders: true,
        includeStyles: true,
        sheetName: 'Sheet1',
        encoding: 'utf-8'
      })
    } else if (newVal === 'pdf') {
      Object.assign(pdfOptions, {
        paperSize: 'A4',
        orientation: 'portrait',
        marginTop: 10,
        marginBottom: 10,
        marginLeft: 15,
        marginRight: 15,
        includePageNumbers: false,
        quality: 'medium'
      })
    } else if (newVal === 'word') {
      Object.assign(wordOptions, {
        includeHeaders: true,
        includeStyles: true,
        template: 'standard'
      })
    } else if (newVal === 'image') {
      Object.assign(imageOptions, {
        format: 'png',
        resolution: 96,
        quality: 90,
        background: 'transparent'
      })
    }
  })

  function handleClose() {
    emit('update:visible', false)
  }

  function handleExport() {
    emit('export', format.value, options.value)
    exporting.value = true
    setTimeout(() => {
      exporting.value = false
      emit('update:visible', false)
    }, 1000)
  }
</script>

<style scoped>
.option-icon {
  margin-right: 8px;
}
</style>
