<template>
  <div class="report-preview">
    <el-card class="preview-toolbar" shadow="never">
      <div class="toolbar-left">
        <span class="toolbar-title">报表预览</span>
        <el-tag size="small" effect="dark" type="info">ID: {{ reportId || '未指定' }}</el-tag>
      </div>
      <div class="toolbar-actions">
        <el-button @click="handleRefresh" :loading="loading">刷新</el-button>
        <el-button @click="printDialog.visible = true">打印预览</el-button>
        <el-button @click="handleExport">导出</el-button>
      </div>
    </el-card>
 
    <ErrorState
       v-if="loadError"
       type="api"
       title="加载预览失败"
       :message="errorMessage"
       :canRetry="true"
       :onRetry="handleRefresh"
    />
 
    <div v-else class="preview-content">
       <LoadingState :loading="loading" text="加载预览..." />
       <div v-if="!renderedHTML" class="preview-empty">
         <el-empty description="暂无预览内容" />
       </div>
       <div v-else class="preview-wrapper">
         <div class="preview-html" v-html="renderedHTML"></div>
         
         <div v-if="showPagination" class="pagination-controls">
           <el-button 
             :disabled="currentPage <= 1" 
             @click="handlePrevPage"
             size="small"
           >
             上一页
           </el-button>
           
           <span class="page-indicator">
             第 {{ currentPage }} / {{ totalPages }} 页
           </span>
           
           <el-button 
             :disabled="currentPage >= totalPages" 
             @click="handleNextPage"
             size="small"
           >
             下一页
           </el-button>
           
           <el-select 
             v-model="pageSize" 
             @change="handlePageSizeChange"
             size="small"
             style="width: 120px; margin-left: 10px;"
           >
             <el-option label="25行/页" :value="25"></el-option>
             <el-option label="50行/页" :value="50"></el-option>
             <el-option label="100行/页" :value="100"></el-option>
             <el-option label="200行/页" :value="200"></el-option>
           </el-select>
         </div>
       </div>
    </div>

    <el-dialog
      v-model="printDialog.visible"
      title="打印预览"
      width="600px"
      :before-close="handlePrintDialogClose"
    >
      <el-form :model="printDialog.settings" label-width="120px">
        <el-form-item label="纸张大小">
          <el-select v-model="printDialog.settings.paperSize">
            <el-option label="A4 (210mm × 297mm)" value="A4" />
            <el-option label="A3 (297mm × 420mm)" value="A3" />
            <el-option label="Letter (216mm × 279mm)" value="Letter" />
          </el-select>
        </el-form-item>

        <el-form-item label="页面方向">
          <el-radio-group v-model="printDialog.settings.orientation">
            <el-radio label="portrait">纵向</el-radio>
            <el-radio label="landscape">横向</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="页边距 (mm)">
          <el-row :gutter="10">
            <el-col :span="12">
              <el-form-item label="上边距" label-width="70px">
                <el-input-number v-model="printDialog.settings.marginTop" :min="0" :max="50" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="下边距" label-width="70px">
                <el-input-number v-model="printDialog.settings.marginBottom" :min="0" :max="50" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="10">
            <el-col :span="12">
              <el-form-item label="左边距" label-width="70px">
                <el-input-number v-model="printDialog.settings.marginLeft" :min="0" :max="50" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="右边距" label-width="70px">
                <el-input-number v-model="printDialog.settings.marginRight" :min="0" :max="50" />
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
              <div class="print-preview-content" v-html="renderedHTML"></div>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="printDialog.visible = false">关闭</el-button>
        <el-button type="primary" @click="handlePrint">打印</el-button>
      </template>
    </el-dialog>
  </div>
</template>
 
<script setup lang="ts">
  import { computed, onMounted, reactive, ref, watch } from 'vue'
  import { useRoute } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { ArrowDown, Setting } from '@element-plus/icons-vue'
  import { reportApi } from '@/api/report'
  import LoadingState from '@/components/common/LoadingState.vue'
  import ErrorState from '@/components/common/ErrorState.vue'
  import PrintPreview from '@/components/report/PrintPreview.vue'
  import ExportDialog from '@/components/report/ExportDialog.vue'

  const route = useRoute()
  const renderedHTML = ref('')
  const loading = ref(false)
  const loadError = ref(false)
  const errorMessage = ref('')
  const currentPage = ref(1)
  const pageSize = ref(50)
  const totalRows = ref(0)
  const exporting = ref(false)
  const reportParams = ref([] as any[])

  const paramDialog = reactive({
     visible: false,
     params: [] as any[],
     values: {} as Record<string, any>
  })

  const exportDialog = reactive({
     visible: false
  })

  const exportProgress = reactive({
     visible: false,
     jobId: ''
  })

  const printDialog = reactive({
     visible: false,
     settings: {
       paperSize: 'A4',
       orientation: 'portrait' as 'portrait' | 'landscape',
       marginTop: 10,
       marginBottom: 10,
       marginLeft: 15,
       marginRight: 15
     }
  })

  const paperSizeDimensions = {
     A4: { width: 210, height: 297 },
     A3: { width: 297, height: 420 },
     Letter: { width: 216, height: 279 }
  }

  const reportId = computed(() => {
   const fromQuery = route.query.id as string | undefined
   if (fromQuery) {
     localStorage.setItem('lastReportId', fromQuery)
     return fromQuery
   }
   return localStorage.getItem('lastReportId') || undefined
  })

  const showPagination = computed(() => {
   return totalRows.value > 50
  })

  const totalPages = computed(() => {
   if (!showPagination.value || pageSize.value === 0) {
     return 1
   }
   return Math.ceil(totalRows.value / pageSize.value)
  })

  const printPreviewStyle = computed(() => {
    const dim = paperSizeDimensions[printDialog.settings.paperSize as keyof typeof paperSizeDimensions]
    const isLandscape = printDialog.settings.orientation === 'landscape'

    return {
      width: `${dim.height}mm`,
      height: `${dim.width}mm`,
      padding: `${printDialog.settings.marginTop}mm ${printDialog.settings.marginRight}mm ${printDialog.settings.marginBottom}mm ${printDialog.settings.marginLeft}mm`,
      background: '#f5f5f5',
      border: '1px solid #ddd',
      transform: isLandscape ? 'rotate(-90deg)' : 'none',
      transformOrigin: 'center center',
      overflow: 'hidden'
    }
  })

  watch(printDialog.settings, () => {
    applyPrintStyles()
  }, { deep: true })

  async function handleRefresh() {
    if (!reportId.value) {
      ElMessage.warning('未指定报表 ID，请先在报表设计器中保存报表')
      return
    }

    loadError.value = false
    errorMessage.value = ''
    loading.value = true
    try {
      const response = await reportApi.preview({ id: reportId.value, params: {} })
      if (response.data.success) {
        renderedHTML.value = response.data.result?.html || ''
        if (response.data.result?.totalRows !== undefined) {
          totalRows.value = response.data.result.totalRows
        }
      } else {
        loadError.value = true
        errorMessage.value = response.data.message || '加载报表失败'
      }
    } catch (error: any) {
      loadError.value = true
      errorMessage.value = error.message || '加载报表失败'
    } finally {
      loading.value = false
    }
  }

  function handleParamDialogClose() {
    if (reportParams.value.length > 0) {
      ElMessage.warning('请输入必填参数')
    } else {
      paramDialog.visible = false
    }
  }

  async function handlePreviewWithParams() {
    const missingParams: string[] = []
    reportParams.value.forEach(p => {
      if (p.required && (paramDialog.values[p.code] === undefined || paramDialog.values[p.code] === '')) {
        missingParams.push(p.name)
      }
    })

    if (missingParams.length > 0) {
      ElMessage.warning(`请填写必填参数: ${missingParams.join(', ')}`)
      return
    }

    paramDialog.visible = false
    currentPage.value = 1
    await loadPreviewData(paramDialog.values)
  }

  async function handlePrevPage() {
    if (currentPage.value > 1) {
      currentPage.value--
      await loadPreviewData({})
    }
  }

  async function handleNextPage() {
    if (currentPage.value < totalPages.value) {
      currentPage.value++
      await loadPreviewData({})
    }
  }

  async function handlePageSizeChange(newSize: number) {
    currentPage.value = 1
    pageSize.value = newSize
    await loadPreviewData({})
  }

  function handleExport() {
    exportDialog.visible = true
  }

  async function handleExportWithOptions(format: string, exportOptions: any) {
    if (!reportId.value) {
      ElMessage.warning('未指定报表 ID，请先在报表设计器中保存报表')
      return
    }

    exportDialog.visible = false
    exporting.value = true
    try {
      const params = {
        id: reportId.value,
        format: format,
        options: exportOptions
      }

      const response = await reportApi.export(params) as any

      if (response.data.success && response.data.result && response.data.result.jobId) {
        exportProgress.jobId = response.data.result.jobId
        exportProgress.visible = true
        ElMessage.success('导出任务已创建')
      } else {
        throw new Error(response.data.message || '导出失败')
      }
    } catch (error: any) {
      ElMessage.error(error.message || '导出失败')
    } finally {
      exporting.value = false
    }
  }

  function handlePrintPreview() {
    if (!renderedHTML.value) {
      ElMessage.warning('请先加载报表预览')
      return
    }
    printDialog.visible = true
  }

  function handlePrintDialogClose() {
    removePrintStyles()
    printDialog.visible = false
  }

  function applyPrintStyles() {
    removePrintStyles()

    const style = document.createElement('style')
    style.id = 'print-styles'
    const s = printDialog.settings
    const isLandscape = s.orientation === 'landscape'

    let css = '@page { '
    css += `size: ${s.paperSize} ${s.orientation}; `
    css += `margin: ${s.marginTop}mm ${s.marginRight}mm ${s.marginBottom}mm ${s.marginLeft}mm; `
    css += '} '

    css += '@media print { '
    css += 'body * { visibility: hidden; } '
    css += '#print-area, #print-area * { visibility: visible; } '
    css += 'body { margin: 0; } '
    css += '}'

    style.textContent = css
    document.head.appendChild(style)

    const previewArea = document.querySelector('.preview-html') as HTMLElement
    if (previewArea) {
      previewArea.id = 'print-area'
    }
  }

  function removePrintStyles() {
    const style = document.getElementById('print-styles')
    if (style) {
      style.remove()
    }

    const previewArea = document.querySelector('.preview-html') as HTMLElement
    if (previewArea && previewArea.id === 'print-area') {
      previewArea.removeAttribute('id')
    }
  }

  function handlePrint() {
    applyPrintStyles()

    setTimeout(() => {
      window.print()
    }, 100)
  }

  async function loadPreviewData(params: Record<string, any>) {
    if (!reportId.value) return
    
    loading.value = true
    try {
      const previewParams: any = {
        id: reportId.value,
        params: params || {}
      }

      if (showPagination.value && currentPage.value > 0) {
        previewParams.params.page = currentPage.value
        previewParams.params.pageSize = pageSize.value
      }

      const response = await reportApi.preview(previewParams)
      if (response.data.success) {
        renderedHTML.value = response.data.result?.html || ''
      } else {
        loadError.value = true
        errorMessage.value = response.data.message || '预览加载失败'
      }
    } catch (error: any) {
      loadError.value = true
      errorMessage.value = error.message || '预览加载失败'
    } finally {
      loading.value = false
    }
  }

  onMounted(async () => {
    if (reportId.value) {
      await handleRefresh()
    }
  })
</script>

<style scoped>
.report-preview {
  min-height: 100vh;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: radial-gradient(circle at top left, #eef2f6, #f8fafc 35%, #ffffff);
  font-family: "IBM Plex Sans", "Noto Sans", sans-serif;
}

.preview-toolbar {
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  background: linear-gradient(145deg, #0f172a, #1f2937);
  color: #f8fafc;
}

.preview-toolbar :deep(.el-card__body) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 16px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-title {
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 0.6px;
}

.toolbar-actions {
  display: flex;
  gap: 10px;
}

.preview-content {
  flex: 1;
  border-radius: 14px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  box-shadow: 0 12px 24px rgba(15, 23, 42, 0.08);
  padding: 16px;
  overflow: auto;
}

.preview-wrapper {
  display: flex;
  flex-direction: column;
}

.preview-empty {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-html {
  min-height: 400px;
}

.preview-html :deep(table) {
  border-collapse: collapse;
}

.preview-html :deep(td),
.preview-html :deep(th) {
  border: 1px solid #e2e8f0;
  padding: 6px 10px;
}

.pagination-controls {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 16px 0;
  border-top: 1px solid #e2e8f0;
  margin-top: 16px;
}

.page-indicator {
  font-size: 14px;
  color: #64748b;
}

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
