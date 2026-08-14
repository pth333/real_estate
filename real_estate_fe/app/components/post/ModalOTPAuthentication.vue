<template>
    <n-modal v-if="show" v-model:show="show" :style="{ maxWidth: '440px' }" :mask-closable="phoneVerified" preset="card"
        :title="step === 'phone' ? 'Xác thực số điện thoại' : 'Nhập mã OTP'" :closable="phoneVerified">

        <div class="flex flex-col gap-5">
            <!-- Error message -->
            <n-alert v-if="errorText" type="error" closable @close="errorText = ''" class="text-sm!">
                {{ errorText }}
            </n-alert>

            <!-- Bước 1: Nhập số điện thoại -->
            <template v-if="step === 'phone'">
                <p class="text-sm text-gray-500">Vui lòng nhập số điện thoại để xác thực trước khi đăng tin</p>

                <n-input v-model:value="phone" placeholder="Nhập số điện thoại" size="large" maxlength="11">
                    <template #prefix>
                        <n-icon size="20">
                            <IconPhone />
                        </n-icon>
                    </template>
                </n-input>

                <n-button type="primary" size="large" :loading="sending" :disabled="!phone.trim()"
                    @click="handleSendOTP">
                    Gửi mã OTP
                </n-button>
            </template>

            <!-- Bước 2: Nhập OTP -->
            <template v-else>
                <p class="text-sm text-gray-500">
                    Mã OTP đã được gửi đến <strong>{{ phone }}</strong>
                </p>

                <n-input v-model:value="otp" placeholder="Nhập mã OTP (6 số)" size="large" maxlength="6">
                    <template #prefix>
                        <n-icon size="20">
                            <IconLock />
                        </n-icon>
                    </template>
                </n-input>

                <div class="flex gap-3 justify-end">
                    <n-button size="large" @click="goBack">
                        Quay lại
                    </n-button>
                    <n-button type="primary" size="large" :loading="verifying" :disabled="otp.length < 6"
                        @click="handleVerifyOTP">
                        Xác thực
                    </n-button>
                </div>

                <p class="text-xs text-center text-gray-400">
                    <span v-if="countdown > 0">Gửi lại mã sau {{ countdown }}s</span>
                    <n-button v-else text size="tiny" @click="handleResendOTP">
                        Gửi lại mã OTP
                    </n-button>
                </p>
            </template>
        </div>
    </n-modal>
</template>

<script setup lang="ts">
interface SendOTPResponse {
    success: boolean
    message?: string
}

interface VerifyOTPResponse {
    success: boolean
    message?: string
}

const emit = defineEmits<{
    close: []
}>()

const { $api } = useNuxtApp()
const { setPhoneVerified, phoneVerified } = usePhoneVerification()

const show = defineModel<boolean>('show', { default: false })

const step = ref<'phone' | 'otp'>('phone')
const phone = ref('')
const otp = ref('')

const sending = ref(false)
const verifying = ref(false)
const countdown = ref(0)
const errorText = ref('')
let countdownTimer: ReturnType<typeof setInterval> | null = null

const handleSendOTP = async () => {
    if (!phone.value.trim()) return

    sending.value = true
    try {
        const res = await $api.post<SendOTPResponse>('/auth/send-otp', {
            phone: phone.value.trim()
        })

        if (res.success) {
            errorText.value = ''
            step.value = 'otp'
            startCountdown()
        } else {
            errorText.value = res.message || 'Gửi OTP thất bại'
        }
    } catch {
        errorText.value = 'Không thể gửi OTP, vui lòng thử lại'
    } finally {
        sending.value = false
    }
}

const handleVerifyOTP = async () => {
    if (otp.value.length < 6) return

    verifying.value = true
    try {
        const res = await $api.post<VerifyOTPResponse>('/auth/verify-otp', {
            phone: phone.value.trim(),
            otp: otp.value
        })

        if (res.success) {
            if (countdownTimer) clearInterval(countdownTimer)
            errorText.value = ''
            setPhoneVerified(phone.value.trim())
            show.value = false
            // Reset state
            step.value = 'phone'
            phone.value = ''
            otp.value = ''
        } else {
            errorText.value = res.message || 'Mã OTP không đúng'
        }
    } catch {
        errorText.value = 'Xác thực thất bại, vui lòng thử lại'
    } finally {
        verifying.value = false
    }
}

const handleResendOTP = () => {
    otp.value = ''
    handleSendOTP()
}

const goBack = () => {
    step.value = 'phone'
    otp.value = ''
    if (countdownTimer) clearInterval(countdownTimer)
    countdown.value = 0
}

const startCountdown = () => {
    countdown.value = 120
    if (countdownTimer) clearInterval(countdownTimer)
    countdownTimer = setInterval(() => {
        countdown.value--
        if (countdown.value <= 0) {
            if (countdownTimer) clearInterval(countdownTimer)
        }
    }, 1000)
}

onUnmounted(() => {
    if (countdownTimer) clearInterval(countdownTimer)
})
</script>
