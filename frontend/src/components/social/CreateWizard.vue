<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="visible"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs"
        @click.self="emit('close')"
        role="dialog"
        aria-modal="true"
        :aria-label="`Create ${entityType}`"
      >
        <div
          class="mx-4 w-full max-w-md rounded-2xl bg-surface-raised p-6 shadow-2xl ring-1 ring-white/10
                 motion-safe:animate-modal-in"
        >
          <!-- Step indicator -->
          <div class="flex items-center gap-2 mb-6">
            <div
              v-for="step in 3"
              :key="step"
              class="h-1 flex-1 rounded-full transition-colors duration-300"
              :class="step <= createStep ? 'bg-spotify' : 'bg-white/10'"
            />
          </div>

          <!-- Step 1: Name -->
          <template v-if="createStep === 1">
            <h2 class="text-lg font-bold text-white">
              {{ entityType === 'party' ? 'Start a Party' : entityType === 'room' ? 'Go Live' : 'Create Club' }}
            </h2>
            <p class="mt-1 text-sm text-white/40">Choose a name for your {{ entityType }}</p>
            <input
              v-model="form.name"
              type="text"
              :placeholder="entityType === 'club' ? 'Club name' : 'Party name'"
              aria-label="Name"
              class="mt-4 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder-white/20 outline-hidden transition focus:border-white/20 focus:bg-white/8"
              dir="auto"
            />
          </template>

          <!-- Step 2: Details -->
          <template v-if="createStep === 2">
            <h2 class="text-lg font-bold text-white">Details</h2>
            <p class="mt-1 text-sm text-white/40">Add a description and set privacy</p>
            <textarea
              v-model="form.description"
              placeholder="What's this about?"
              rows="3"
              aria-label="Description"
              class="mt-4 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder-white/20 outline-hidden transition focus:border-white/20 focus:bg-white/8"
              dir="auto"
            />
            <label class="mt-4 flex items-center gap-3 cursor-pointer rounded-xl bg-white/4 p-3 transition hover:bg-white/6">
              <input
                v-model="form.isPublic"
                type="checkbox"
                class="h-5 w-5 rounded border-white/10 bg-white/5 accent-[#1db954]"
              />
              <div>
                <span class="text-sm text-white">Public</span>
                <p class="text-xs text-white/30">Anyone can find and join</p>
              </div>
            </label>
          </template>

          <!-- Step 3: Review -->
          <template v-if="createStep === 3">
            <div class="flex flex-col items-center gap-4 py-4 text-center">
              <div class="flex h-16 w-16 items-center justify-center rounded-full bg-spotify/10">
                <i aria-hidden="true" :class="entityIcon" class="text-2xl text-spotify" />
              </div>
              <h2 class="text-lg font-bold text-white">Almost there!</h2>
              <p class="text-sm text-white/40">Review and launch your {{ entityType }}</p>
              <div class="w-full rounded-xl bg-white/4 p-4 text-start">
                <p class="text-xs text-white/30">Name</p>
                <p class="text-sm font-medium text-white">{{ form.name }}</p>
                <p v-if="form.description" class="mt-3 text-xs text-white/30">Description</p>
                <p v-if="form.description" class="text-sm text-white/60">{{ form.description }}</p>
                <p class="mt-3 text-xs text-white/30">Visibility</p>
                <p class="text-sm text-white/60">
                  <i aria-hidden="true" :class="form.isPublic ? 'pi pi-globe' : 'pi pi-lock'" class="text-xs me-1" />
                  {{ form.isPublic ? 'Public' : 'Private' }}
                </p>
              </div>
            </div>
          </template>

          <!-- Navigation buttons -->
          <div class="mt-6 flex gap-3">
            <!-- Back / Cancel -->
            <template v-if="createStep === 1">
              <button
                class="flex-1 rounded-xl bg-white/5 py-3 text-sm font-medium text-white/50 transition hover:bg-white/10"
                @click="emit('close')"
              >
                Cancel
              </button>
            </template>
            <template v-else>
              <button
                class="flex-1 rounded-xl bg-white/5 py-3 text-sm font-medium text-white/50 transition hover:bg-white/10"
                @click="createStep--"
              >
                Back
              </button>
            </template>

            <!-- Next / Create -->
            <template v-if="createStep < 3">
              <button
                :disabled="createStep === 1 && !form.name.trim()"
                class="flex-1 rounded-xl bg-spotify py-3 text-sm font-bold text-black transition hover:bg-spotify-hover disabled:opacity-40 focus-visible:outline-2 focus-visible:outline-white"
                @click="createStep++"
              >
                Next
              </button>
            </template>
            <template v-else>
              <button
                :disabled="creating"
                class="flex-1 rounded-xl bg-spotify py-3 text-sm font-bold text-black transition hover:bg-spotify-hover disabled:opacity-40 focus-visible:outline-2 focus-visible:outline-white"
                @click="handleCreate"
              >
                <span v-if="creating" class="inline-flex items-center gap-2">
                  <span class="inline-block h-4 w-4 animate-spin rounded-full border-2 border-black border-t-transparent" />
                  Creating...
                </span>
                <span v-else>
                  <i aria-hidden="true" class="pi pi-send text-xs me-1.5" />
                  Launch
                </span>
              </button>
            </template>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useSocialApi } from '@/services/api/social'

const props = defineProps<{
  visible: boolean
  entityType: 'party' | 'room' | 'club'
}>()

const emit = defineEmits<{
  close: []
  created: [id: string, type: string]
}>()

const api = useSocialApi()
export type WizardEntityType = 'party' | 'room' | 'club'

const createStep = ref(1)
const creating = ref(false)

const form = reactive({
  name: '',
  description: '',
  isPublic: true,
})

const entityIcon = computed(() => {
  switch (props.entityType) {
    case 'party': return 'pi pi-users'
    case 'room': return 'pi pi-megaphone'
    case 'club': return 'pi pi-building'
  }
})

function reset() {
  createStep.value = 1
  form.name = ''
  form.description = ''
  form.isPublic = true
  creating.value = false
}

async function handleCreate() {
  if (!form.name.trim() || creating.value) return
  creating.value = true
  try {
    let result: { id?: string } | undefined

    if (props.entityType === 'party') {
      result = await api.createParty({
        title: form.name.trim(),
        description: form.description.trim() || undefined,
        is_public: form.isPublic,
      })
    } else if (props.entityType === 'room') {
      result = await api.createRoom({
        title: form.name.trim(),
        description: form.description.trim() || undefined,
        is_public: form.isPublic,
      })
    } else {
      const slug = form.name.trim()
        .toLowerCase()
        .replace(/[^a-z0-9\u0600-\u06FF\s-]/g, '')
        .replace(/\s+/g, '-')
        .replace(/-+/g, '-')
        .replace(/^-|-$/g, '')
      result = await api.createClub({
        name: form.name.trim(),
        slug,
        description: form.description.trim() || undefined,
        is_public: form.isPublic,
      })
    }

    if (result?.id) {
      emit('created', result.id, props.entityType)
    }
    emit('close')
  } catch (err) {
    console.error('Failed to create:', err)
  } finally {
    creating.value = false
  }
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
@keyframes modal-in {
  from { opacity: 0; transform: scale(0.96) translateY(8px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}
.animate-modal-in {
  animation: modal-in 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}
</style>
