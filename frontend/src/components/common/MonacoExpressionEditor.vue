<template>
  <div ref="editorContainer" class="monaco-expression-editor"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as monaco from 'monaco-editor'

interface Field {
  name: string
  displayName: string
  type?: string
  dataType?: string
}

interface FunctionDef {
  name: string
  signature: string
  description?: string
  category: string
}

interface Props {
  modelValue: string
  fields: Field[]
  functions?: FunctionDef[]
  height?: string
  readOnly?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  height: '150px',
  readOnly: false,
  functions: () => []
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
  (e: 'blur'): void
  (e: 'focus'): void
}>()

const editorContainer = ref<HTMLElement>()
let editor: monaco.editor.IStandaloneCodeEditor | null = null

// Register expression language
const registerExpressionLanguage = () => {
  // Register the expression language
  monaco.languages.register({ id: 'goreport-expression' })

  // Define syntax highlighting rules
  monaco.languages.setMonarchTokensProvider('goreport-expression', {
    defaultToken: '',
    tokenPostfix: '.expr',

    keywords: [
      'IF', 'CASE', 'WHEN', 'THEN', 'ELSE', 'END',
      'AND', 'OR', 'NOT', 'IN', 'IS', 'NULL'
    ],

    operators: [
      '+', '-', '*', '/', '%', '=', '!=', '<>', '<', '>', '<=', '>=',
      '||', '&&', '!'
    ],

    functions: [
      'SUM', 'AVG', 'COUNT', 'MAX', 'MIN',
      'CONCAT', 'SUBSTRING', 'LENGTH', 'UPPER', 'LOWER', 'TRIM',
      'DATE_FORMAT', 'DATE_ADD', 'DATE_SUB', 'DATEDIFF', 'NOW', 'CURDATE', 'CURTIME',
      'YEAR', 'MONTH', 'DAY', 'HOUR', 'MINUTE', 'SECOND',
      'ROUND', 'CEIL', 'FLOOR', 'ABS'
    ],

    symbols: /[=><!~?:&|+\-*\/\^%]+/,

    tokenizer: {
      root: [
        // Field references like [fieldName]
        [/\[[\w\s]+\]/, 'field.reference'],

        // Numbers
        [/\d*\.\d+([eE][\-+]?\d+)?/, 'number.float'],
        [/\d+/, 'number'],

        // Strings
        [/'([^'\\]|\\.)*'/, 'string'],
        [/"([^"\\]|\\.)*"/, 'string'],

        // Functions
        [/[a-zA-Z_]\w*(?=\()/, {
          cases: {
            '@functions': 'function',
            '@keywords': 'keyword',
            '@default': 'identifier'
          }
        }],

        // Keywords
        [/[a-zA-Z_]\w*/, {
          cases: {
            '@keywords': 'keyword',
            '@default': 'identifier'
          }
        }],

        // Operators
        [/@symbols/, {
          cases: {
            '@operators': 'operator',
            '@default': ''
          }
        }],

        // Whitespace
        [/\s+/, 'white'],

        // Delimiters
        [/[()]/, 'delimiter.parenthesis'],
        [/[\[\]]/, 'delimiter.bracket'],
        [/[{},]/, 'delimiter'],
      ]
    }
  })

  // Configure language configuration (brackets, auto-indent, etc.)
  monaco.languages.setLanguageConfiguration('goreport-expression', {
    brackets: [
      ['(', ')'],
      ['[', ']']
    ],
    autoClosingPairs: [
      { open: '(', close: ')' },
      { open: '[', close: ']' },
      { open: "'", close: "'" },
      { open: '"', close: '"' }
    ],
    surroundingPairs: [
      { open: '(', close: ')' },
      { open: '[', close: ']' },
      { open: "'", close: "'" },
      { open: '"', close: '"' }
    ]
  })
}

// Register completion provider
const registerCompletionProvider = () => {
  monaco.languages.registerCompletionItemProvider('goreport-expression', {
    triggerCharacters: ['[', '(', ',', ' '],

    provideCompletionItems: (model, position) => {
      const word = model.getWordUntilPosition(position)
      const range = {
        startLineNumber: position.lineNumber,
        endLineNumber: position.lineNumber,
        startColumn: word.startColumn,
        endColumn: word.endColumn
      }

      const suggestions: monaco.languages.CompletionItem[] = []

      // Add field suggestions when typing [ or at start
      const lineContent = model.getLineContent(position.lineNumber)
      const charBeforeCursor = lineContent[position.column - 2]

      if (charBeforeCursor === '[' || word.word.startsWith('[')) {
        props.fields.forEach(field => {
          suggestions.push({
            label: field.displayName || field.name,
            kind: monaco.languages.CompletionItemKind.Field,
            documentation: `字段: ${field.name}\n类型: ${field.dataType || 'unknown'}`,
            insertText: field.name,
            range,
            detail: field.type || 'field'
          })
        })
      }

      // Add function suggestions
      if (props.functions && props.functions.length > 0) {
        props.functions.forEach(fn => {
          suggestions.push({
            label: fn.name,
            kind: monaco.languages.CompletionItemKind.Function,
            documentation: fn.description || `${fn.category}函数: ${fn.signature}`,
            insertText: fn.signature,
            insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            range,
            detail: fn.category
          })
        })
      }

      // Add keyword suggestions
      const keywords = ['IF', 'CASE', 'WHEN', 'THEN', 'ELSE', 'END', 'AND', 'OR', 'NOT']
      keywords.forEach(keyword => {
        suggestions.push({
          label: keyword,
          kind: monaco.languages.CompletionItemKind.Keyword,
          insertText: keyword,
          range
        })
      })

      return { suggestions }
    }
  })
}

// Register hover provider for tooltips
const registerHoverProvider = () => {
  monaco.languages.registerHoverProvider('goreport-expression', {
    provideHover: (model, position) => {
      const word = model.getWordAtPosition(position)
      if (!word) return null

      const wordText = word.word.toUpperCase()

      // Check if it's a function
      const func = props.functions?.find(f => f.name.toUpperCase() === wordText)
      if (func) {
        return {
          contents: [
            { value: `**${func.name}** - ${func.category}` },
            { value: `\`\`\`\n${func.signature}\n\`\`\`` },
            { value: func.description || '' }
          ]
        }
      }

      // Check if it's a field
      const field = props.fields.find(f => f.name.toUpperCase() === wordText)
      if (field) {
        return {
          contents: [
            { value: `**${field.displayName || field.name}**` },
            { value: `字段名: \`${field.name}\`\n类型: ${field.dataType || 'unknown'}\n类别: ${field.type || 'field'}` }
          ]
        }
      }

      return null
    }
  })
}

// Initialize editor
const initEditor = () => {
  if (!editorContainer.value) return

  // Register language features
  registerExpressionLanguage()
  registerCompletionProvider()
  registerHoverProvider()

  // Create editor
  editor = monaco.editor.create(editorContainer.value, {
    value: props.modelValue,
    language: 'goreport-expression',
    theme: 'vs',
    minimap: { enabled: false },
    lineNumbers: 'off',
    folding: false,
    lineDecorationsWidth: 0,
    lineNumbersMinChars: 0,
    glyphMargin: false,
    scrollbar: {
      vertical: 'auto',
      horizontal: 'auto'
    },
    overviewRulerLanes: 0,
    hideCursorInOverviewRuler: true,
    renderLineHighlight: 'none',
    quickSuggestions: true,
    suggestOnTriggerCharacters: true,
    wordBasedSuggestions: 'off',
    parameterHints: { enabled: true },
    formatOnType: true,
    autoIndent: 'advanced',
    readOnly: props.readOnly,
    fontSize: 14,
    fontFamily: 'Consolas, "Courier New", monospace',
    padding: { top: 8, bottom: 8 }
  })

  // Handle content changes
  editor.onDidChangeModelContent(() => {
    const value = editor?.getValue() || ''
    emit('update:modelValue', value)
    emit('change', value)
  })

  // Handle focus/blur
  editor.onDidBlurEditorWidget(() => {
    emit('blur')
  })

  editor.onDidFocusEditorWidget(() => {
    emit('focus')
  })
}

// Watch for external value changes
watch(() => props.modelValue, (newValue) => {
  if (editor && editor.getValue() !== newValue) {
    editor.setValue(newValue)
  }
})

// Watch for readOnly changes
watch(() => props.readOnly, (newValue) => {
  if (editor) {
    editor.updateOptions({ readOnly: newValue })
  }
})

// Watch for fields changes (update completion provider)
watch(() => props.fields, () => {
  // Re-register completion provider with new fields
  registerCompletionProvider()
}, { deep: true })

onMounted(() => {
  initEditor()
})

onUnmounted(() => {
  if (editor) {
    editor.dispose()
    editor = null
  }
})

// Public methods
const insertText = (text: string) => {
  if (!editor) return

  const position = editor.getPosition()
  if (position) {
    editor.executeEdits('insert', [
      {
        range: new monaco.Range(
          position.lineNumber,
          position.column,
          position.lineNumber,
          position.column
        ),
        text
      }
    ])
    editor.focus()
  }
}

const getEditorInstance = () => editor

defineExpose({
  insertText,
  getEditorInstance
})
</script>

<style scoped>
.monaco-expression-editor {
  width: 100%;
  height: v-bind(height);
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

.monaco-expression-editor :deep(.monaco-editor) {
  padding-top: 4px;
}
</style>
