export const useHeroBg = () => {
  const HERO_BG_URL = 'https://pub-5eb4e976c2fe4062ba3cdabce48568cc.r2.dev/uploads/931c3b1968a2e54e32e00dace86c38de.jpg'

  const backgroundStyle = computed(() => ({
    backgroundImage: `url('${HERO_BG_URL}')`
  }))

  return { HERO_BG_URL, backgroundStyle }
}