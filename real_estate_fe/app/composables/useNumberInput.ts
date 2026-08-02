/**
 * Composable xử lý ô nhập số có định dạng dấu chấm hàng nghìn (vi-VN).
 *
 * Nhận vào một Ref<number | null> trỏ thẳng vào field của store.
 * Trả về `displayed` (Ref<string>) để bind vào v-model:value của n-input.
 * Trả về `onInput` để bind vào @update:value — format ngay khi gõ và ghi số thực về store.
 *
 * Ví dụ:
 *   const { displayed, onInput } = useNumberInput(toRef(postStore.form, 'price'))
 *   <n-input :value="displayed" @update:value="onInput" />
 */
export function useNumberInput(sourceRef: Ref<number | null>) {
    // Khởi tạo từ giá trị hiện có trong store
    const displayed = ref<string>(
        sourceRef.value !== null ? sourceRef.value.toLocaleString('vi-VN') : '',
    )

    // Khi người dùng gõ: bỏ ký tự không phải số, format, ghi về store
    const onInput = (val: string) => {
        const digits = val.replace(/\D/g, '')
        displayed.value = digits ? Number(digits).toLocaleString('vi-VN') : ''
        sourceRef.value = digits ? Number(digits) : null
    }

    // Đồng bộ khi store thay đổi từ bên ngoài (reset form...)
    watch(sourceRef, (newVal) => {
        const currentNum = Number(displayed.value.replace(/\./g, '')) || null
        if (currentNum !== newVal) {
            displayed.value = newVal !== null ? newVal.toLocaleString('vi-VN') : ''
        }
    })

    return { displayed, onInput }
}
