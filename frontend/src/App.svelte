<script lang="ts">
	import { api } from "./lib/api";

	type Role = "select" | "librarian" | "student";

	type Group = {
		code: string;
		specialty_code: string;
		admission_year: number;
		institute_no: number;
		department_no: number;
		group_no: number;
		current_course: number;
	};

	type Student = {
		id: number;
		firstname: string;
		lastname: string;
	};

	type Book = {
		id: number;
		name: string;
		authors: string[];
		disciplineid: number;
		course: number;
		isbn: string;
		publicationyear: number;
		publisher: string;
		available_count?: number;
		total_count?: number;
		discipline?: string;
	};

	type CopyList = {
		copy_ids: number[];
	};

	type Loan = {
		id: number;
		book_id: number;
		book_name: string;
		authors: string[];
		copy_id: number;
		issue_date: string;
	};

	type IssueRequest = {
		student_id: number;
		copy_id: number;
	};

	type IssueResponse = {
		id: number;
	};

	let role: Role = "select";

	let groupQuery = "";
	let bookQuery = "";

	let group: Group | null = null;
	let students: Student[] = [];
	let recommendations: Book[] = [];
	let searchResults: Book[] = [];
	let studentLoans: Loan[] = [];

	let selectedStudentId: number | null = null;
	let selectedBookId: number | null = null;
	let selectedBook: Book | null = null;
	let selectedCopyId: number | null = null;
	let availableCopyIds: number[] = [];

	let loadingGroup = false;
	let loadingBook = false;
	let loadingSearch = false;
	let saving = false;

	let error = "";
	let notice = "";

	$: currentStudent =
		selectedStudentId !== null
			? (students.find((s) => s.id === selectedStudentId) ?? null)
			: null;

	$: currentCourse = group?.current_course ?? null;
	$: currentAuthors = selectedBook?.authors?.length
		? selectedBook.authors.join(", ")
		: "—";

	function resetFeedback() {
		error = "";
		notice = "";
	}

	function toArray<T>(value: T[] | null | undefined): T[] {
		return Array.isArray(value) ? value : [];
	}

	function authorsText(book: Book) {
		return Array.isArray(book.authors) && book.authors.length > 0
			? book.authors.join(", ")
			: "—";
	}

	function currentBookCount(book: Book) {
		return book.available_count ?? 0;
	}

	function todayIso() {
		const d = new Date();
		const yyyy = d.getFullYear();
		const mm = String(d.getMonth() + 1).padStart(2, "0");
		const dd = String(d.getDate()).padStart(2, "0");
		return `${yyyy}-${mm}-${dd}`;
	}

	async function loadSelectedBook(id: number) {
		loadingBook = true;
		try {
			const [book, copies] = await Promise.all([
				api.get<Book>(`/api/books/${id}`),
				api.get<CopyList>(`/api/books/${id}/copies`),
			]);

			selectedBook = book ?? null;
			selectedBookId = id;
			availableCopyIds = toArray(copies?.copy_ids);
			selectedCopyId = availableCopyIds[0] ?? null;
		} finally {
			loadingBook = false;
		}
	}

	async function loadStudentLoans(studentId: number) {
		studentLoans = toArray(
			await api.get<Loan[]>(`/api/students/${studentId}/loans`),
		);
	}

	async function loadGroup() {
		const code = groupQuery.trim();
		if (!code) {
			error = "Введите код группы";
			return;
		}

		loadingGroup = true;
		resetFeedback();

		try {
			const [g, s, r] = await Promise.all([
				api.get<Group>(`/api/groups/${code}`),
				api.get<Student[]>(`/api/groups/${code}/students`),
				api.get<Book[]>(`/api/groups/${code}/recommendations`),
			]);

			group = g ?? null;
			students = toArray(s);
			recommendations = toArray(r);

			if (students.length > 0) {
				if (
					selectedStudentId === null ||
					!students.some((x) => x.id === selectedStudentId)
				) {
					selectedStudentId = students[0].id;
				}

				await loadStudentLoans(selectedStudentId);
			} else {
				selectedStudentId = null;
				studentLoans = [];
			}

			if (recommendations.length > 0) {
				await loadSelectedBook(recommendations[0].id);
			} else {
				selectedBook = null;
				selectedBookId = null;
				selectedCopyId = null;
				availableCopyIds = [];
			}

			notice = `Группа ${code} загружена`;
		} catch (e) {
			error = e instanceof Error ? e.message : "Ошибка загрузки группы";
		} finally {
			loadingGroup = false;
		}
	}

	async function searchBooks() {
		const q = bookQuery.trim();
		if (!q) {
			searchResults = [];
			return;
		}

		loadingSearch = true;
		resetFeedback();

		try {
			searchResults = toArray(
				await api.get<Book[]>(
					`/api/books/search?q=${encodeURIComponent(q)}`,
				),
			);
		} catch (e) {
			error = e instanceof Error ? e.message : "Ошибка поиска книг";
		} finally {
			loadingSearch = false;
		}
	}

	async function selectStudent(studentId: number) {
		selectedStudentId = studentId;
		resetFeedback();

		try {
			studentLoans = toArray(
				await api.get<Loan[]>(`/api/students/${studentId}/loans`),
			);
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Ошибка загрузки выдач студента";
		}
	}

	async function selectBook(id: number) {
		resetFeedback();

		try {
			await loadSelectedBook(id);
		} catch (e) {
			error = e instanceof Error ? e.message : "Ошибка загрузки книги";
		}
	}

	function patchBookState(
		bookId: number,
		copyId: number | null,
		delta: number,
	) {
		const patch = (book: Book) =>
			book.id === bookId
				? {
						...book,
						available_count: Math.max(
							0,
							(book.available_count ?? 0) + delta,
						),
					}
				: book;

		recommendations = recommendations.map(patch);
		searchResults = searchResults.map(patch);

		if (selectedBook?.id === bookId) {
			selectedBook = patch(selectedBook);
		}

		if (selectedBookId === bookId) {
			if (delta < 0 && copyId !== null) {
				availableCopyIds = availableCopyIds.filter(
					(id) => id !== copyId,
				);
				selectedCopyId = availableCopyIds[0] ?? null;
			}

			if (delta > 0 && copyId !== null) {
				if (!availableCopyIds.includes(copyId)) {
					availableCopyIds = [...availableCopyIds, copyId].sort(
						(a, b) => a - b,
					);
				}
				if (selectedCopyId === null) {
					selectedCopyId = availableCopyIds[0] ?? copyId;
				}
			}
		}
	}

	async function issueBook() {
		if (
			selectedStudentId === null ||
			selectedCopyId === null ||
			selectedBookId === null ||
			selectedBook === null
		) {
			error = "Выберите студента, книгу и экземпляр";
			return;
		}

		saving = true;
		resetFeedback();

		const payload: IssueRequest = {
			student_id: selectedStudentId,
			copy_id: selectedCopyId,
		};

		try {
			const created = await api.post<IssueResponse>(
				"/api/issuing",
				payload,
			);

			patchBookState(selectedBookId, selectedCopyId, -1);

			studentLoans = [
				{
					id: created.id,
					book_id: selectedBookId,
					book_name: selectedBook.name,
					authors: selectedBook.authors ?? [],
					copy_id: selectedCopyId,
					issue_date: todayIso(),
				},
				...studentLoans,
			];

			notice = "Книга выдана";
		} catch (e) {
			error = e instanceof Error ? e.message : "Ошибка выдачи книги";
		} finally {
			saving = false;
		}
	}

	async function returnBook(loan: Loan) {
		saving = true;
		resetFeedback();

		try {
			await api.patch(`/api/issuing/${loan.id}/return`);

			patchBookState(loan.book_id, loan.copy_id, 1);
			studentLoans = studentLoans.filter((x) => x.id !== loan.id);

			notice = "Книга возвращена";
		} catch (e) {
			error = e instanceof Error ? e.message : "Ошибка возврата книги";
		} finally {
			saving = false;
		}
	}
</script>

<svelte:head>
	<title>Библиотека ГУАП</title>
</svelte:head>

{#if role === "select"}
	<div class="page">
		<div class="shell">
			<div class="topbar">
				<div>
					<div class="title">Библиотека ГУАП</div>
					<div class="subtitle">
						Выбор должности и вход в рабочий режим
					</div>
				</div>
				<div class="badges">
					<span>Библиотека</span>
					<span>ГУАП</span>
					<span>UI</span>
				</div>
			</div>

			<div class="landing">
				<div class="card hero">
					<div>
						<div class="hero-title">Интерфейс библиотечной БД</div>
						<div class="hero-text">
							Подбор литературы по группе и студенту, поиск книг,
							выдача и возврат экземпляров.
						</div>
					</div>
					<div class="chips">
						<span>Поиск</span>
						<span>Рекомендации</span>
						<span>Выдача</span>
					</div>
				</div>

				<div class="card">
					<div class="section-title">Выберите должность</div>

					<button
						class="role-btn"
						on:click={() => (role = "librarian")}
					>
						<div class="role-index blue">01</div>
						<div>
							<div class="role-name">Библиотекарь</div>
							<div class="role-desc">
								Поиск студента, рекомендации, выдача и возврат
							</div>
						</div>
					</button>

					<button
						class="role-btn"
						on:click={() => (role = "student")}
					>
						<div class="role-index violet">02</div>
						<div>
							<div class="role-name">Студент</div>
							<div class="role-desc">
								Поиск книг и просмотр рекомендаций
							</div>
						</div>
					</button>

					<button
						class="role-btn"
						on:click={() => (role = "librarian")}
					>
						<div class="role-index red">03</div>
						<div>
							<div class="role-name">Администратор</div>
							<div class="role-desc">Технический режим</div>
						</div>
					</button>
				</div>
			</div>
		</div>
	</div>
{:else if role === "librarian"}
	<div class="page">
		<div class="shell">
			<div class="topbar">
				<div>
					<div class="title">Рабочий кабинет библиотекаря</div>
					<div class="subtitle">
						Поиск группы, рекомендации и операции с выдачей
					</div>
				</div>
				<div class="badges">
					<span>ГУАП</span>
					<span>Библиотека</span>
					<button
						class="small-link"
						on:click={() => (role = "select")}
					>
						Сменить роль
					</button>
				</div>
			</div>

			<div class="dashboard">
				<div class="col">
					<div class="card">
						<div class="section-title">Группа</div>

						<div class="input-row">
							<input
								name="searchGroup"
								bind:value={groupQuery}
								placeholder="7322"
							/>
							<button
								class="action blue"
								on:click={loadGroup}
								disabled={loadingGroup}
							>
								{loadingGroup ? "Загрузка..." : "Найти"}
							</button>
						</div>

						<div class="mini-grid">
							<div class="field">
								<div class="label">Специальность</div>
								<div class="value">
									{group?.specialty_code ?? "—"}
								</div>
							</div>
							<div class="field">
								<div class="label">Курс</div>
								<div class="value">{currentCourse ?? "—"}</div>
							</div>
							<div class="field">
								<div class="label">Институт</div>
								<div class="value">
									{group?.institute_no ?? "—"}
								</div>
							</div>
							<div class="field">
								<div class="label">Кафедра</div>
								<div class="value">
									{group?.department_no ?? "—"}
								</div>
							</div>
						</div>
					</div>

					<div class="card">
						<div class="section-title">Студенты группы</div>

						<div class="stack">
							{#each students as student}
								<button
									class:selected={student.id ===
										selectedStudentId}
									class="list-row"
									on:click={() => selectStudent(student.id)}
								>
									<div>
										<div class="row-title">
											{student.lastname}
											{student.firstname}
										</div>
										<div class="row-sub">
											Группа {group?.code ?? "—"}
										</div>
									</div>
									<span class="tag">Студент</span>
								</button>
							{/each}

							{#if students.length === 0}
								<div class="empty">Список студентов пуст.</div>
							{/if}
						</div>
					</div>

					<div class="card">
						<div class="section-title">Выдачи студента</div>

						{#if currentStudent}
							<div class="hint">
								{currentStudent.lastname}
								{currentStudent.firstname}
							</div>

							<div class="stack">
								{#each studentLoans as loan}
									<div class="list-row static">
										<div>
											<div class="row-title">
												{loan.book_name}
											</div>
											<div class="row-sub">
												Экземпляр {loan.copy_id} · выдано
												{loan.issue_date}
											</div>
										</div>

										<button
											class="tag red"
											type="button"
											on:click={() => returnBook(loan)}
											disabled={saving}
										>
											Возврат
										</button>
									</div>
								{/each}

								{#if studentLoans.length === 0}
									<div class="empty">
										У выбранного студента нет активных
										выдач.
									</div>
								{/if}
							</div>
						{:else}
							<div class="empty">
								Сначала выберите студента группы.
							</div>
						{/if}

						{#if error}
							<div class="message error">{error}</div>
						{/if}
					</div>
				</div>

				<div class="col">
					<div class="card">
						<div class="section-title">Оформление выдачи</div>

						<div class="mini-grid">
							<div class="field">
								<div class="label">Студент</div>
								<div class="value">
									{currentStudent
										? `${currentStudent.lastname} ${currentStudent.firstname}`
										: "—"}
								</div>
							</div>
							<div class="field">
								<div class="label">Книга</div>
								<div class="value">
									{selectedBook?.name ?? "—"}
								</div>
							</div>
							<div class="field">
								<div class="label">Экземпляр</div>
								<div class="value">{selectedCopyId ?? "—"}</div>
							</div>
							<div class="field">
								<div class="label">Статус</div>
								<div class="value">
									{availableCopyIds.length > 0
										? "готов к выдаче"
										: "нет доступных экземпляров"}
								</div>
							</div>
						</div>

						<div class="actions">
							<button
								class="action blue"
								on:click={issueBook}
								disabled={saving ||
									selectedStudentId === null ||
									selectedCopyId === null}
							>
								{saving ? "Сохранение..." : "Выдать"}
							</button>
						</div>
					</div>

					<div class="card">
						<div class="section-title">
							Рекомендованная литература
						</div>
						<div class="hint">
							По группе {group?.code ?? "—"} и специальности {group?.specialty_code ??
								"—"}
						</div>

						<div class="stack">
							{#each recommendations as book}
								<button
									class:selected={book.id === selectedBookId}
									class="book-row"
									on:click={() => selectBook(book.id)}
								>
									<div>
										<div class="row-title">{book.name}</div>
										<div class="row-sub">
											{authorsText(book)} · {book.isbn}
										</div>
									</div>
									<span class="tag">
										{currentBookCount(book)} экз.
									</span>
								</button>
							{/each}

							{#if recommendations.length === 0}
								<div class="empty">
									Рекомендации не найдены.
								</div>
							{/if}
						</div>
					</div>

					<div class="card">
						<div class="section-title">Поиск книги</div>

						<div class="input-row">
							<input
								bind:value={bookQuery}
								placeholder="Базы данных"
							/>
							<button
								class="action blue"
								on:click={searchBooks}
								disabled={loadingSearch}
							>
								{loadingSearch ? "Поиск..." : "Поиск"}
							</button>
						</div>

						<div class="stack">
							{#each searchResults as book}
								<button
									class="book-row"
									on:click={() => selectBook(book.id)}
								>
									<div>
										<div class="row-title">{book.name}</div>
										<div class="row-sub">
											{authorsText(book)} · курс {book.course}
										</div>
									</div>
									<span class="tag">
										{currentBookCount(book)} экз.
									</span>
								</button>
							{/each}

							{#if bookQuery.trim() && searchResults.length === 0}
								<div class="empty">
									По запросу ничего не найдено.
								</div>
							{/if}
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
{:else}
	<div class="page">
		<div class="shell">
			<div class="topbar">
				<div>
					<div class="title">Студенческий режим</div>
					<div class="subtitle">
						Поиск книг и просмотр рекомендаций
					</div>
				</div>
				<div class="badges">
					<span>Студент</span>
					<button
						class="small-link"
						on:click={() => (role = "select")}
					>
						Сменить роль
					</button>
				</div>
			</div>

			<div class="dashboard student-layout">
				<div class="col">
					<div class="card">
						<div class="section-title">Группа</div>

						<div class="input-row">
							<input bind:value={groupQuery} placeholder="8226" />
							<button
								class="action blue"
								on:click={loadGroup}
								disabled={loadingGroup}
							>
								{loadingGroup ? "Загрузка..." : "Найти"}
							</button>
						</div>

						<div class="mini-grid">
							<div class="field">
								<div class="label">Специальность</div>
								<div class="value">
									{group?.specialty_code ?? "—"}
								</div>
							</div>
							<div class="field">
								<div class="label">Курс</div>
								<div class="value">{currentCourse ?? "—"}</div>
							</div>
						</div>
					</div>

					<div class="card">
						<div class="section-title">Рекомендации</div>

						<div class="stack">
							{#each recommendations as book}
								<button
									class:selected={book.id === selectedBookId}
									class="book-row"
									on:click={() => selectBook(book.id)}
								>
									<div>
										<div class="row-title">{book.name}</div>
										<div class="row-sub">
											{authorsText(book)} · курс {book.course}
										</div>
									</div>
									<span class="tag">
										{currentBookCount(book)} экз.
									</span>
								</button>
							{/each}

							{#if recommendations.length === 0}
								<div class="empty">
									Для группы ничего не найдено.
								</div>
							{/if}
						</div>
					</div>
				</div>

				<div class="col">
					<div class="card">
						<div class="section-title">Поиск книги</div>

						<div class="input-row">
							<input
								bind:value={bookQuery}
								placeholder="Базы данных"
							/>
							<button class="action blue" on:click={searchBooks}>
								Поиск
							</button>
						</div>

						<div class="stack">
							{#each searchResults as book}
								<button
									class:selected={book.id === selectedBookId}
									class="book-row"
									on:click={() => selectBook(book.id)}
								>
									<div>
										<div class="row-title">{book.name}</div>
										<div class="row-sub">
											{authorsText(book)} · {book.publisher}
										</div>
									</div>
									<span class="tag">
										{currentBookCount(book)} экз.
									</span>
								</button>
							{/each}

							{#if bookQuery.trim() && searchResults.length === 0}
								<div class="empty">
									По запросу ничего не найдено.
								</div>
							{/if}
						</div>
					</div>

					<div class="card">
						<div class="section-title">Карточка книги</div>

						<div class="book-card">
							<div class="book-title">
								{selectedBook?.name ?? "—"}
							</div>
							<div class="book-author">{currentAuthors}</div>

							<div class="mini-grid">
								<div class="field">
									<div class="label">ISBN</div>
									<div class="value">
										{selectedBook?.isbn ?? "—"}
									</div>
								</div>
								<div class="field">
									<div class="label">Год</div>
									<div class="value">
										{selectedBook?.publicationyear ?? "—"}
									</div>
								</div>
								<div class="field">
									<div class="label">Издательство</div>
									<div class="value">
										{selectedBook?.publisher ?? "—"}
									</div>
								</div>
								<div class="field">
									<div class="label">Курс</div>
									<div class="value">
										{selectedBook?.course ?? "—"}
									</div>
								</div>
							</div>
						</div>
					</div>

					{#if error}
						<div class="message error">{error}</div>
					{/if}
					{#if notice}
						<div class="message ok">{notice}</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	:global(body) {
		margin: 0;
		font-family: Arial, sans-serif;
		background: #eef3f9;
		color: #1f2d3d;
	}

	.page {
		min-height: 100vh;
		padding: 16px;
		box-sizing: border-box;
	}

	.shell {
		max-width: 1400px;
		margin: 0 auto;
		background: #fff;
		border: 1px solid #d7e0ea;
		border-radius: 20px;
		overflow: hidden;
		box-shadow: 0 12px 30px rgba(31, 45, 61, 0.06);
	}

	.topbar {
		background: #005aaa;
		color: #fff;
		padding: 18px 20px;
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 16px;
	}

	.title {
		font-size: 22px;
		font-weight: 700;
		line-height: 1.1;
	}

	.subtitle {
		margin-top: 4px;
		font-size: 13px;
		color: rgba(255, 255, 255, 0.8);
	}

	.badges {
		display: flex;
		gap: 8px;
		align-items: center;
		flex-wrap: wrap;
		justify-content: flex-end;
	}

	.badges span,
	.small-link {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 6px 10px;
		border-radius: 999px;
		border: none;
		font-size: 12px;
		font-weight: 700;
	}

	.badges span:nth-child(1) {
		background: #edf4ff;
		color: #005aaa;
	}

	.badges span:nth-child(2) {
		background: #f7eef6;
		color: #ab3a8d;
	}

	.badges span:nth-child(3) {
		background: #fdecef;
		color: #e4003a;
	}

	.small-link {
		background: transparent;
		color: #fff;
		border: 1px solid rgba(255, 255, 255, 0.3);
		cursor: pointer;
	}

	.landing,
	.dashboard {
		padding: 18px;
		background: #f7faff;
	}

	.landing {
		display: grid;
		grid-template-columns: 1.2fr 0.8fr;
		gap: 14px;
	}

	.dashboard {
		display: grid;
		grid-template-columns: 0.95fr 1.25fr;
		gap: 14px;
	}

	.student-layout {
		grid-template-columns: 1fr 1.1fr;
	}

	.col {
		display: flex;
		flex-direction: column;
		gap: 14px;
	}

	.card {
		background: #fff;
		border: 1px solid #d7e0ea;
		border-radius: 18px;
		padding: 16px;
		box-shadow: 0 8px 18px rgba(31, 45, 61, 0.04);
	}

	.hero {
		min-height: 240px;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
	}

	.hero-title {
		font-size: 28px;
		font-weight: 700;
		margin-bottom: 8px;
	}

	.hero-text {
		font-size: 15px;
		line-height: 1.5;
		color: #5c6d7e;
		max-width: 520px;
	}

	.chips {
		display: flex;
		gap: 8px;
		flex-wrap: wrap;
		margin-top: 18px;
	}

	.chips span,
	.tag {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 6px 10px;
		border-radius: 999px;
		font-size: 12px;
		font-weight: 700;
		border: none;
	}

	.chips span:nth-child(1),
	.blue {
		background: #edf4ff;
		color: #005aaa;
	}

	.chips span:nth-child(2),
	.violet {
		background: #f7eef6;
		color: #ab3a8d;
	}

	.chips span:nth-child(3),
	.red {
		background: #fdecef;
		color: #e4003a;
	}

	.section-title {
		font-size: 15px;
		font-weight: 700;
		margin-bottom: 12px;
	}

	.role-btn,
	.list-row,
	.book-row {
		width: 100%;
		border: 1px solid #d7e0ea;
		background: #f7faff;
		border-radius: 14px;
		padding: 12px;
		text-align: left;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
		cursor: pointer;
		box-sizing: border-box;
	}

	.role-btn:hover,
	.list-row:hover,
	.book-row:hover {
		border-color: #005aaa;
	}

	.role-index {
		width: 38px;
		height: 38px;
		border-radius: 999px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 700;
		flex: 0 0 auto;
	}

	.role-name {
		font-weight: 700;
		margin-bottom: 2px;
	}

	.role-desc {
		font-size: 12px;
		color: #5c6d7e;
	}

	.input-row {
		display: flex;
		gap: 8px;
		margin-bottom: 12px;
	}

	input {
		flex: 1;
		border: 1px solid #d7e0ea;
		border-radius: 12px;
		padding: 11px 12px;
		font-size: 14px;
		background: #fff;
	}

	.action {
		border: none;
		border-radius: 12px;
		padding: 11px 14px;
		font-weight: 700;
		cursor: pointer;
		color: #fff;
	}

	.action.blue {
		background: #005aaa;
	}

	.action:disabled {
		opacity: 0.7;
		cursor: not-allowed;
	}

	.mini-grid {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: 10px;
	}

	.field {
		border: 1px solid #d7e0ea;
		background: #f7faff;
		border-radius: 12px;
		padding: 10px;
	}

	.label {
		font-size: 12px;
		color: #8595a7;
		margin-bottom: 4px;
	}

	.value {
		font-size: 14px;
		font-weight: 700;
	}

	.stack {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.row-title {
		font-weight: 700;
		margin-bottom: 2px;
	}

	.row-sub {
		font-size: 12px;
		color: #5c6d7e;
	}

	.list-row.selected,
	.book-row.selected {
		border-color: #005aaa;
		background: #edf4ff;
	}

	.list-row.static {
		cursor: default;
	}

	.empty,
	.hint {
		font-size: 12px;
		color: #8595a7;
	}

	.actions {
		margin-top: 12px;
		display: flex;
		justify-content: flex-end;
	}

	.book-card {
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.book-title {
		font-size: 22px;
		font-weight: 700;
	}

	.book-author {
		color: #5c6d7e;
	}

	.message {
		border-radius: 12px;
		padding: 10px 12px;
		font-size: 13px;
	}

	.message.error {
		background: #fdecef;
		color: #8d1f36;
	}

	.message.ok {
		background: #edf4ff;
		color: #005aaa;
	}

	.tag.red {
		cursor: pointer;
	}

	@media (max-width: 1100px) {
		.landing,
		.dashboard,
		.student-layout {
			grid-template-columns: 1fr;
		}
	}
</style>
