import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "klient KSeF",
  description: "Dokumentacja użytkownika",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Dokumentacja', link: '/content/' }
    ],

    sidebar: [
      {
        text: 'Instalacja',
        items: [
          { text: 'Instalacja programu', link: '/content/instalacja' },
        ]
      },
      {
        text: 'Komendy',
        items: [
          { text: 'Zapisanie tokenu', link: '/content/komendy/save-token'},
          { text: 'Generowanie faktur', link: '/content/komendy/generate'},
          { text: 'Wysyłka faktur', items: [
            {text: 'Sesja wsadowa (batch)', link: '/content/komendy/upload/batch'},
            {text: 'Sesja interaktywna', link: '/content/komendy/upload/interaktywna'},
            {text: 'Kody QR', link: '/content/qrcode'},
          ]},
          { text: 'Pobieranie faktur', link: '/content/komendy/download'},
          { text: 'Pobieranie UPO', link: '/content/komendy/upo'},
          { text: 'Identyfikator płatności', link: '/content/komendy/payment-id'},
          { text: 'Wizualizacja PDF', link: '/content/komendy/wizualizacja-pdf'},
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/toudi/ksef' }
    ]
  }
})
