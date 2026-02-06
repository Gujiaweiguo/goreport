<template>
  <div class="chart-property-panel">
    <div class="panel-header">
      <h3>图表属性</h3>
      <el-tag size="small" effect="plain" type="warning">Live</el-tag>
    </div>

    <el-form ref="formRef" :model="formState" :rules="rules" label-position="top" class="property-form">
      <el-collapse v-model="activeGroups">
        <el-collapse-item title="基础属性" name="basic">
          <el-form-item label="标题" prop="title">
            <el-input v-model="formState.title" placeholder="请输入图表标题" @input="emitChange" />
          </el-form-item>
          <div class="grid-row">
            <el-form-item label="宽度" prop="width">
              <el-input-number v-model="formState.width" :min="280" :max="1600" :step="10" @change="emitChange" />
            </el-form-item>
            <el-form-item label="高度" prop="height">
              <el-input-number v-model="formState.height" :min="220" :max="1000" :step="10" @change="emitChange" />
            </el-form-item>
          </div>
          <el-form-item label="边距" prop="margin">
            <el-input-number v-model="formState.margin" :min="0" :max="80" :step="1" @change="emitChange" />
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="样式属性" name="style">
          <el-form-item label="主题" prop="theme">
            <el-select v-model="formState.theme" @change="emitChange">
              <el-option label="默认" value="default" />
              <el-option label="深色" value="dark" />
            </el-select>
          </el-form-item>
          <el-form-item label="主色" prop="color">
            <el-color-picker v-model="formState.color" show-alpha @change="emitChange" />
          </el-form-item>
          <div class="grid-row">
            <el-form-item label="字体" prop="fontFamily">
              <el-select v-model="formState.fontFamily" @change="emitChange">
                <el-option label="Noto Serif SC" value="'Noto Serif SC', serif" />
                <el-option label="Source Han Sans" value="'Source Han Sans CN', sans-serif" />
                <el-option label="IBM Plex Sans" value="'IBM Plex Sans', sans-serif" />
              </el-select>
            </el-form-item>
            <el-form-item label="字号" prop="fontSize">
              <el-input-number v-model="formState.fontSize" :min="10" :max="40" :step="1" @change="emitChange" />
            </el-form-item>
          </div>
          <el-form-item label="图例" prop="showLegend">
            <el-switch v-model="formState.showLegend" @change="emitChange" />
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="交互属性" name="interaction">
          <el-form-item label="鼠标悬停提示" prop="hoverable">
            <el-switch v-model="formState.hoverable" @change="emitChange" />
          </el-form-item>
          <el-form-item label="点击交互" prop="clickable">
            <el-switch v-model="formState.clickable" @change="emitChange" />
          </el-form-item>
          <el-form-item label="缩放" prop="zoomable">
            <el-switch v-model="formState.zoomable" @change="emitChange" />
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="标题配置" name="title-config">
          <el-form-item label="主标题" prop="mainTitle">
            <el-input v-model="formState.mainTitle" placeholder="请输入主标题" @input="emitChange" />
          </el-form-item>
          <el-form-item label="副标题" prop="subTitle">
            <el-input v-model="formState.subTitle" placeholder="请输入副标题" @input="emitChange" />
          </el-form-item>
          <el-form-item label="标题位置" prop="titlePosition">
            <el-select v-model="formState.titlePosition" @change="emitChange">
              <el-option label="左侧" value="left" />
              <el-option label="居中" value="center" />
              <el-option label="右侧" value="right" />
            </el-select>
          </el-form-item>
        </el-collapse-item>
      </el-collapse>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'

export interface ChartPropertyConfig {
  title: string
  width: number
  height: number
  margin: number
  theme: 'default' | 'dark'
  color: string
  fontFamily: string
  fontSize: number
  showLegend: boolean
  hoverable: boolean
  clickable: boolean
  zoomable: boolean
  mainTitle: string
  subTitle: string
  titlePosition: 'left' | 'center' | 'right'
}

const props = withDefaults(
  defineProps<{
    modelValue: ChartPropertyConfig
  }>(),
  {
    modelValue: () => ({
      title: '图表标题',
      width: 760,
      height: 420,
      margin: 16,
      theme: 'default',
      color: '#2d6cdf',
      fontFamily: "'Noto Serif SC', serif",
      fontSize: 14,
      showLegend: true,
      hoverable: true,
      clickable: true,
      zoomable: false,
      mainTitle: '图表标题',
      subTitle: '副标题',
      titlePosition: 'left'
    })
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: ChartPropertyConfig]
  change: [value: ChartPropertyConfig]
}>()

const formRef = ref<FormInstance>()
const activeGroups = ref(['basic', 'style', 'interaction', 'title-config'])

const formState = reactive<ChartPropertyConfig>({
  title: '图表标题',
  width: 760,
  height: 420,
  margin: 16,
  theme: 'default',
  color: '#2d6cdf',
  fontFamily: "'Noto Serif SC', serif",
  fontSize: 14,
  showLegend: true,
  hoverable: true,
  clickable: true,
  zoomable: false,
  mainTitle: '图表标题',
  subTitle: '副标题',
  titlePosition: 'left'
})

const rules = reactive<FormRules<ChartPropertyConfig>>({
  title: [{ required: true, message: '标题不能为空', trigger: 'blur' }],
  width: [{ required: true, message: '宽度不能为空', trigger: 'change' }],
  height: [{ required: true, message: '高度不能为空', trigger: 'change' }],
  margin: [{ required: true, message: '边距不能为空', trigger: 'change' }],
  fontSize: [{ required: true, message: '字号不能为空', trigger: 'change' }]
})

function syncForm(value: ChartPropertyConfig) {
  formState.title = value.title
  formState.width = value.width
  formState.height = value.height
  formState.margin = value.margin
  formState.theme = value.theme
  formState.color = value.color
  formState.fontFamily = value.fontFamily
  formState.fontSize = value.fontSize
  formState.showLegend = value.showLegend
  formState.hoverable = value.hoverable
  formState.clickable = value.clickable
  formState.zoomable = value.zoomable
  formState.mainTitle = value.mainTitle
  formState.subTitle = value.subTitle
  formState.titlePosition = value.titlePosition
}

function emitChange() {
  const nextValue: ChartPropertyConfig = {
    title: formState.title,
    width: formState.width,
    height: formState.height,
    margin: formState.margin,
    theme: formState.theme,
    color: formState.color,
    fontFamily: formState.fontFamily,
    fontSize: formState.fontSize,
    showLegend: formState.showLegend,
    hoverable: formState.hoverable,
    clickable: formState.clickable,
    zoomable: formState.zoomable,
    mainTitle: formState.mainTitle,
    subTitle: formState.subTitle,
    titlePosition: formState.titlePosition
  }
  emit('update:modelValue', nextValue)
  emit('change', nextValue)
  void formRef.value?.validate().catch(() => {
    return undefined
  })
}

watch(
  () => props.modelValue,
  value => {
    syncForm(value)
  },
  { immediate: true, deep: true }
)
</script>

<style scoped>
.chart-property-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  border: 1px solid #dfe6f2;
  background: #ffffff;
  overflow: hidden;
}

.panel-header {
  height: 46px;
  padding: 0 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e8edf5;
  background: linear-gradient(180deg, #f8faff 0%, #f2f6fc 100%);
}

.panel-header h3 {
  margin: 0;
  font-size: 14px;
  color: #1f2937;
}

.property-form {
  flex: 1;
  overflow: auto;
  padding: 10px 12px;
}

.grid-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

:deep(.el-form-item) {
  margin-bottom: 10px;
}

:deep(.el-input-number),
:deep(.el-select),
:deep(.el-input) {
  width: 100%;
}

:deep(.el-collapse) {
  border-top: none;
  border-bottom: none;
}

:deep(.el-collapse-item__header) {
  color: #334155;
  font-weight: 600;
}
</style>
