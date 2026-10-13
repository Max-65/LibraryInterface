#set page(width: 297mm, height: auto, margin: 8mm, fill: rgb("#eef3f9"))
#set text(font: ("Roboto", "Arial"), size: 9pt, fill: rgb("#1f2d3d"))

#let blue = rgb("#005aaa")
#let violet = rgb("#ab3a8d")
#let red = rgb("#e4003a")
#let surface = rgb("#ffffff")
#let soft = rgb("#f7faff")
#let line = rgb("#d7e0ea")
#let muted = rgb("#5c6d7e")
#let faint = rgb("#8595a7")
#let pale_blue = rgb("#edf4ff")
#let pale_violet = rgb("#f7eef6")
#let pale_red = rgb("#fdecef")

#let pill(label, fill, ink) = rect(
  fill: fill,
  radius: 999pt,
  inset: (x: 8pt, y: 4pt),
)[#text(fill: ink, size: 8pt, weight: "bold")[#label]]

#let shell(title, subtitle, body, inset: 12pt) = rect(
  fill: surface,
  radius: 16pt,
  inset: 0pt,
  stroke: 1pt + line,
)[
  #rect(fill: blue, radius: (top-left: 16pt, top-right: 16pt, bottom-left: 0pt, bottom-right: 0pt), inset: 12pt)[
    #grid(
      columns: (1fr, auto),
      gutter: 10pt,
      [
        #text(fill: surface, size: 16pt, weight: "bold")[#title]
        #v(0pt)
        #text(fill: rgb("#d7e8ff"), size: 8.5pt)[#subtitle]
      ],
    )
  ]
  #v(weak: true, 0pt)
  #rect(fill: soft, inset: inset)[#body]
]

#let card(title, body, accent: blue, height: auto, width: auto) = rect(
  fill: surface,
  radius: 14pt,
  inset: 12pt,
  stroke: 1pt + line,
  height: height,
  width: width,
)[
  #rect(fill: accent, radius: 999pt, inset: (x: 8pt, y: 3pt))[
    #text(fill: surface, size: 8pt, weight: "bold")[#title]
  ]
  #v(0pt)
  #body
]

#let field(label, value, height: auto, width: auto) = rect(
  fill: surface,
  radius: 12pt,
  inset: 10pt,
  stroke: 1pt + line,
  height: height,
  width: width,
)[
  #text(fill: faint, size: 7.5pt)[#label]
  #v(0pt)
  #text(size: 10pt, weight: "medium")[#value]
]

#let small_stat(label, value, accent) = rect(
  fill: surface,
  radius: 14pt,
  inset: 12pt,
  stroke: 1pt + line,
)[
  #text(fill: muted, size: 8pt)[#label]
  #v(6pt)
  #text(fill: accent, size: 18pt, weight: "bold")[#value]
]

#let button(label, fill: blue, radius: 10pt, inset: (x: 10pt, y: 7pt)) = rect(
  fill: fill,
  radius: radius,
  inset: inset,
)[#text(fill: surface, size: 8.5pt, weight: "bold")[#label]]

#let row_item(a, b, c, tagfill: pale_blue, tagink: blue) = rect(
  fill: surface,
  radius: 12pt,
  inset: 10pt,
  stroke: 1pt + line,
)[
  #grid(
    columns: (3fr, 4fr, auto),
    gutter: 8pt,
    [
      #text(size: 10pt, weight: "bold")[#a]
      #v(2pt)
      #text(fill: muted, size: 8pt)[#b]
    ],
    align(left)[
      #text(fill: faint, size: 8pt)[#c]
    ],
    pill("Подробнее", tagfill, tagink),
  )
]

#let book_item(title, authors, meta, status, statusfill, statusink) = rect(
  fill: surface,
  radius: 12pt,
  inset: 10pt,
  stroke: 1pt + line,
)[
  #grid(
    columns: (1fr, auto),
    gutter: 8pt,
    [
      #text(size: 10pt, weight: "bold")[#title]
      #v(2pt)
      #text(fill: muted, size: 8pt)[#authors]
      #v(2pt)
      #text(fill: faint, size: 7.5pt)[#meta]
    ],
    [
      #pill(status, statusfill, statusink)
    ],
  )
]

#shell(
  "Библиотека ГУАП",
  "Выбор должности и вход в рабочий режим",
  [
    #set align(center + horizon)
    #card(
      "Выберите должность",
      [
        #rect(fill: soft, radius: 12pt, inset: 12pt, stroke: 1pt + line)[
          #grid(
            columns: 1fr,
            gutter: 10pt,
            [
              #text(weight: "bold")[Библиотекарь]
              #v(2pt)
              #text(fill: muted, size: 8pt)[Поиск студента, рекомендации, выдача и возврат]
            ],
          )
        ]
        #v(10pt)
        #rect(fill: soft, radius: 12pt, inset: 12pt, stroke: 1pt + line)[
          #grid(
            columns: 1fr,
            gutter: 10pt,
            [
              #text(weight: "bold")[Студент]
              #v(2pt)
              #text(fill: muted, size: 8pt)[Поиск книг и просмотр рекомендаций]
            ],
          )
        ]
        #v(10pt)
        #rect(fill: soft, radius: 12pt, inset: 12pt, stroke: 1pt + line)[
          #grid(
            columns: 1fr,
            gutter: 10pt,
            [
              #text(weight: "bold")[Администратор]
              #v(2pt)
              #text(fill: muted, size: 8pt)[Технический режим]
            ],
          )
        ]
      ],
      accent: blue,
    )
  ],
  inset: (x: 20%, y: 12pt),
)

#pagebreak()

#shell(
  "Рабочий кабинет библиотекаря",
  "Поиск группы, рекомендации и операции с выдачей",
  [
    #grid(
      columns: (1fr, 1.4fr),
      gutter: 8pt,
      [
        #card(
          "Группа",
          [
            #grid(
              columns: (1fr, 1fr),
              gutter: 8pt,
              [#field("Код группы", "8226", width: 100%)],
              [#align(center + horizon)[#button("Найти", fill: blue, inset: 40%)]],
            )
            #v(0pt)

            #grid(
              columns: (1fr, 1fr),
              gutter: 8pt,

              [
                #field("Специальность", "09.03.03", width: 100%)
              ],

              [
                #field("Курс", "2", width: 100%)
              ],
            )
            #v(0pt)

            #grid(
              columns: (1fr, 1fr),
              gutter: 8pt,

              [
                #field("Институт", "8", width: 100%)
              ],

              [
                #field("Кафедра", "2", width: 100%)
              ],
            )
          ],
          accent: blue,
        )

        #v(12pt)

        #card(
          "Студенты группы",
          [
            #row_item(
              "Иван Петров",
              "8226",
              "2 курс",
              tagfill: pale_blue,
              tagink: blue,
            )

            #v(8pt)

            #row_item(
              "Мария Иванова",
              "8226",
              "2 курс",
              tagfill: pale_violet,
              tagink: violet,
            )

            #v(8pt)

            #row_item(
              "Алексей Смирнов",
              "8226",
              "2 курс",
              tagfill: pale_red,
              tagink: red,
            )
          ],
          accent: violet,
        )

        #v(12pt)

        #card(
          "Активные выдачи",
          [
            #row_item(
              "Иван Петров",
              "Базы данных",
              "до 24.05.2026",
              tagfill: pale_blue,
              tagink: blue,
            )

            #v(8pt)

            #row_item(
              "Мария Иванова",
              "Операционные системы",
              "до 18.05.2026",
              tagfill: pale_violet,
              tagink: violet,
            )
          ],
          accent: red,
        )
      ],

      [
        #card(
          "Рекомендованная литература",
          [
            #grid(
              columns: (1fr, 1fr),
              gutter: 8pt,

              [
                #field("Дисциплина", "Базы данных")
              ],

              [
                #field("Фильтр", "Только доступные")
              ],
            )

            #v(10pt)

            #book_item(
              "Базы данных",
              "Кристофер Дейт",
              "ISBN 978-5-8459-0780-7 · Вильямс · 2018",
              "2 экземпляра",
              pale_blue,
              blue,
            )

            #v(8pt)

            #book_item(
              "SQL для профессионалов",
              "Бен-Форта",
              "ISBN 978-5-496-01234-8 · Питер · 2021",
              "1 экземпляр",
              pale_violet,
              violet,
            )

            #v(8pt)

            #book_item(
              "Java. Библиотека профессионала",
              "Кей Хорстманн",
              "ISBN 978-5-8459-1912-5 · Вильямс · 2019",
              "выдана",
              pale_red,
              red,
            )

            #v(12pt)

            #grid(
              columns: (1fr, 1fr, 1fr),
              gutter: 8pt,

              [
                #button("Карточка книги", fill: blue)
              ],

              [
                #button("Выдать", fill: violet)
              ],

              [
                #button("Возврат", fill: red)
              ],
            )
          ],
          accent: blue,
        )

        #v(12pt)

        #card(
          "Оформление выдачи",
          [
            #grid(
              columns: (1fr, 1fr),
              gutter: 8pt,

              [
                #field("Студент", "Иван Петров")
              ],

              [
                #field("Экземпляр", "copyid 3")
              ],
            )

            #v(8pt)

            #grid(
              columns: (1fr, 1fr),
              gutter: 8pt,

              [
                #field("Книга", "Базы данных")
              ],

              [
                #field("Статус", "готов к выдаче")
              ],
            )

            #v(10pt)

            #grid(
              columns: (1fr, auto),
              gutter: 8pt,

              [
                #text(fill: muted, size: 8pt)[Выдача создаёт запись в issuing, статус экземпляра меняется триггером.]
              ],

              [
                #button("Подтвердить", fill: blue)
              ],
            )
          ],
          accent: violet,
        )
      ],
    )
  ],
)

#pagebreak()

#shell(
  "Каталог и карточка книги",
  "Поиск по названию, просмотр карточки и доступных экземпляров",
  [
    #grid(
      columns: (0.95fr, 1.05fr),
      gutter: 12pt,
      [
        #card(
          "Поиск книги",
          [
            #align(center)[#field("Строка поиска", "Базы данных", width: 50%)]
            #v(10pt)
            #book_item(
              "Базы данных",
              "Кристофер Дейт",
              "09.03.03 · курс 2",
              "id 10",
              pale_blue,
              blue,
            )
            #v(8pt)
            #book_item(
              "SQL для профессионалов",
              "Бен-Форта",
              "09.03.03 · курс 2",
              "id 11",
              pale_violet,
              violet,
            )
            #v(8pt)
            #book_item(
              "Проектирование ПО",
              "И. Соммервилл",
              "09.03.04 · курс 3",
              "id 17",
              pale_red,
              red,
            )
          ],
          accent: blue,
          height: 1fr,
        )
      ],
      [
        #card(
          "Карточка книги",
          [
            #text(size: 20pt, weight: "bold")[Базы данных]
            #v(0pt)
            #text(fill: muted)[Кристофер Дейт]
            #v(10pt)
            #grid(columns: (1fr, 1fr), gutter: 8pt)[
              #field("ISBN", "978-5-8459-0780-7")
              #field("Год", "2018")
            ]
            #v(8pt)
            #grid(columns: (1fr, 1fr), gutter: 8pt)[
              #field("Издательство", "Вильямс")
              #field("Курс", "2")
            ]
            #v(8pt)
            #field("Дисциплина", "Базы данных")
            #v(10pt)
            #grid(columns: (1fr, 1fr, 1fr), gutter: 8pt)[
              #pill("3 всего", pale_blue, blue)
              #pill("2 в наличии", pale_violet, violet)
              #pill("1 выдана", pale_red, red)
            ]
            #v(12pt)
            #button("Показать экземпляры", fill: blue)
          ],
          accent: violet,
        )
        #v(12pt)
        #card(
          "Доступные экземпляры",
          [
            #row_item("Экземпляр 1", "copyid 3", "готов к выдаче", tagfill: pale_blue, tagink: blue)
            #v(8pt)
            #row_item("Экземпляр 2", "copyid 8", "готов к выдаче", tagfill: pale_violet, tagink: violet)
          ],
          accent: red,
        )
      ],
    )
  ],
)
