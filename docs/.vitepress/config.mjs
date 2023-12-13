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
          ]},
          { text: 'Pobieranie UPO', link: '/content/komendy/upo'},
          { text: 'Wizualizacja PDF', link: '/content/komendy/wizualizacja-pdf'},
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/toudi/ksef' }
    ]
  }
})
