import { ref, computed } from 'vue'
import type { DrawerMode, DrawerConfig } from './types'

export function useDrawer(drawerId: string, config: DrawerConfig) {
  const visible = ref(false)
  const mode = ref<DrawerMode>('create')
  const currentEntity = ref<any>(null)

  const openCreate = () => {
    currentEntity.value = null
    mode.value = 'create'
    visible.value = true
  }

  const openEdit = (entity: any) => {
    if (config.beforeOpen) {
      currentEntity.value = config.beforeOpen(entity, 'edit')
    } else {
      currentEntity.value = { ...entity }
    }
    mode.value = 'edit'
    visible.value = true
  }

  const openView = (entity: any) => {
    if (config.beforeOpen) {
      currentEntity.value = config.beforeOpen(entity, 'view')
    } else {
      currentEntity.value = { ...entity }
    }
    mode.value = 'view'
    visible.value = true
  }

  const close = () => {
    visible.value = false
    setTimeout(() => {
      currentEntity.value = null
      mode.value = 'create'
    }, 300)
  }

  const isCreateMode = computed(() => mode.value === 'create')
  const isEditMode = computed(() => mode.value === 'edit')
  const isViewMode = computed(() => mode.value === 'view')

  return {
    drawerId,
    visible,
    mode,
    currentEntity,
    isCreateMode,
    isEditMode,
    isViewMode,
    openCreate,
    openEdit,
    openView,
    close
  }
}