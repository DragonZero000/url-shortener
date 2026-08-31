<template>
  <div class="min-h-screen flex items-center justify-center px-4">
    <div class="w-full max-w-lg">
      <!-- Заголовок -->
      <h1 class="text-3xl font-light text-white mb-8 tracking-tight text-center">
        URL Shortener
      </h1>

      <!-- Форма ввода -->
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <input
            v-model="urlInput"
            type="text"
            placeholder="Введите ссылку для сокращения"
            :class="[
              'w-full px-4 py-3 rounded-lg border text-white placeholder-gray-500 outline-none transition-all duration-200',
              error ? 'border-red-500/60 focus:border-red-500' : 'border-neutral-800 focus:border-indigo-500'
            ]"
            :disabled="loading"
          />
        </div>

        <button
          type="submit"
          :disabled="loading || !urlInput.trim()"
          class="w-full py-3 rounded-lg bg-indigo-600 text-white font-medium hover:bg-indigo-700 disabled:opacity-40 disabled:cursor-not-allowed transition-all duration-200"
        >
          {{ loading ? 'Сокращение...' : 'Сократить' }}
        </button>
      </form>

      <!-- Ошибка -->
      <div v-if="error" class="mt-4 p-3 rounded-lg border border-red-500/30 bg-red-500/5">
        <p class="text-sm text-red-400">{{ error }}</p>
      </div>

      <!-- Результат -->
      <div v-if="result" class="mt-6 p-5 rounded-lg border border-neutral-800 bg-neutral-900/50">
        <div class="mb-4">
          <p class="text-xs text-gray-500 uppercase tracking-wider mb-1.5">Оригинал</p>
          <a :href="originalUrl" target="_blank" rel="noopener noreferrer"
             class="text-sm text-indigo-400 hover:text-indigo-300 break-all transition-colors duration-200">
            {{ originalUrl }}
          </a>
        </div>

        <div class="pt-4 border-t border-neutral-800">
          <p class="text-xs text-gray-500 uppercase tracking-wider mb-1.5">Сокращенная ссылка</p>
          <div class="flex items-center gap-2">
            <a :href="result" target="_blank" rel="noopener noreferrer"
               class="flex-1 text-sm text-indigo-400 hover:text-indigo-300 break-all transition-colors duration-200">
              {{ result }}
            </a>
            <button
              @click="copyToClipboard(result)"
              title="Копировать"
              class="p-2 rounded-md border border-neutral-700 hover:border-indigo-500/60 hover:bg-indigo-500/10 transition-all duration-200 shrink-0"
            >
              <svg v-if="!copied" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-gray-400">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
              </svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#4ade80" stroke-width="2">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { shortenUrl, type ShortenResponse } from '../api/client'

const urlInput = ref('')
const loading = ref(false)
const error = ref<string | null>(null)
const result = ref<string | null>(null)
const copied = ref(false)

function isValidUrl(value: string): boolean {
  try {
    const parsed = new URL(value)
    return ['http:', 'https:'].includes(parsed.protocol)
  } catch {
    return false
  }
}

async function handleSubmit() {
  error.value = null
  result.value = null

  if (!urlInput.value.trim()) {
    error.value = 'Поле не может быть пустым'
    return
  }

  if (!isValidUrl(urlInput.value)) {
    error.value = 'Введите корректный URL (например, https://example.com)'
    return
  }

  loading.value = true
  try {
    const response: ShortenResponse = await shortenUrl(urlInput.value)
    result.value = response.url
  } catch (e) {
    if (e instanceof Error) {
      error.value = e.message || 'Произошла ошибка при сокращении ссылки'
    } else {
      error.value = 'Неизвестная ошибка'
    }
  } finally {
    loading.value = false
  }
}

async function copyToClipboard(text: string) {
  try {
    await navigator.clipboard.writeText(text)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    // Fallback для старых браузеров
    const textarea = document.createElement('textarea')
    textarea.value = text
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  }
}

const originalUrl = computed(() => urlInput.value)
</script>