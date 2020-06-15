package incrementor

import (
	"fmt"
	"sync"
	"testing"
)

// Запуск тестов:
// go test -v -run=Test -cover incrementor_test.go incrementor.go
// Результат:
// ok      command-line-arguments  0.212s  coverage: 100.0% of statements

// Запуск бэнчмарков:
// go test -run=Bench -benchmem -bench=. incrementor_test.go incrementor.go
// Результат:
// BenchmarkIncrementor-24         20000000                87.6 ns/op             0 B/op          0 allocs/op

// Тест стандартного использования инкрементора
func TestIncrementorBasicOperations(t *testing.T) {
	// Устанвливаем количество вызово инкремента значения и ожижаемый результат
	testCases := []struct {
		times int
		want  int
	}{
		{0, 0},
		{1, 1},
		{3, 3},
		{42, 42},
		{43, 0},
		{44, 1},
	}

	// Цикл по тестам
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v in %v", tc.times, tc.want), func(t *testing.T) {
			// Инициализируем инкрементор
			var inc Incrementor = NewIncrementor()
			// Устанавливаем максимальное значение равное 42
			if err := inc.SetMaximumValue(42); err != nil {
				t.Fatal(err)
			}
			// Инкрементируем установленное количество раз
			for k := 0; k < tc.times; k++ {
				inc.IncrementNumber()
			}
			// Получаем значение инкрементора
			if got := inc.GetNumber(); got != tc.want {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}
}

// Тест установки максимального значения по умолчанию
func TestIncrementorWithoutSetMaximumValue(t *testing.T) {
	const want = MaximumInt
	t.Run(fmt.Sprintf("%v in %v", MaximumInt, MaximumInt), func(t *testing.T) {
		inc := NewIncrementor()
		inc.IncrementNumber()
		if got := inc.GetNumber(); got != 1 {
			t.Errorf("got %v; want %v", got, want)
		}
	})
}

// Тест появления ошибки при попытке установить отрицательное значение максимума
func TestIncrementorSetMaximumValueError(t *testing.T) {
	t.Run(fmt.Sprintf("%v in %v", MaximumInt, MaximumInt), func(t *testing.T) {
		inc := NewIncrementor()
		err := inc.SetMaximumValue(-42)
		if err == nil || err != ErrNegativeMaximumValue {
			t.Errorf("got %v; want %v", err, ErrNegativeMaximumValue)
		}
	})
}

// Тест обнуления значения при установке максимума меньше актуального значения
func TestIncrementorSetZeroValue(t *testing.T) {
	const want = 0
	t.Run(fmt.Sprintf("%v in %v", 0, 0), func(t *testing.T) {
		inc := NewIncrementor()
		inc.IncrementNumber() //1
		inc.IncrementNumber() //2
		inc.IncrementNumber() //3
		if err := inc.SetMaximumValue(2); err != nil {
			t.Fatal(err)
		}
		if got := inc.GetNumber(); got != want {
			t.Errorf("got %v; want %v", got, want)
		}
	})
}

// Тест дублирования инкрементора
func TestIncrementorTwiceWithIncrementNumber(t *testing.T) {
	const want = 2
	t.Run(fmt.Sprintf("%v in %v", 0, 0), func(t *testing.T) {
		inc1 := NewIncrementor()
		inc1.IncrementNumber() //1
		inc2 := inc1
		inc1.IncrementNumber() //2
		if got := inc2.GetNumber(); got != want {
			t.Errorf("got %v; want %v", got, want)
		}
	})
}

// Тест обнуления значения у дублированного инкрементора
func TestIncrementorTwiceWithSetMaximumValue(t *testing.T) {
	const want = 0
	t.Run(fmt.Sprintf("%v in %v", 0, 0), func(t *testing.T) {
		inc1 := NewIncrementor()
		inc1.IncrementNumber() //1
		inc2 := inc1
		inc1.IncrementNumber() //2
		if err := inc1.SetMaximumValue(1); err != nil {
			t.Fatal(err)
		}
		if got := inc2.GetNumber(); got != want {
			t.Errorf("got %v; want %v", got, want)
		}
	})
}

// Тест конкурентного использования инкрементора
func TestIncrementorConcurrence(t *testing.T) {
	const kCount = 50  // Количество горутин
	const pCount = 100 // Количество вызовов инкремента в каждой горутине
	const want = kCount * pCount
	t.Run(fmt.Sprintf("%v in %v", MaximumInt, MaximumInt), func(t *testing.T) {
		inc := NewIncrementor()
		var wg sync.WaitGroup // Инициализация WaitGroup для ожидания завершения работы всех горутин
		for k := 0; k < kCount; k++ {
			wg.Add(1)
			go func() {
				for p := 0; p < pCount; p++ {
					inc.IncrementNumber()
				}
				wg.Done()
			}()
		}
		wg.Wait()
		if got := inc.GetNumber(); got != want {
			t.Errorf("got %v; want %v", got, want)
		}
	})
}

//// Тест на установку значения превышающего максимум через unsafe
//func TestIncrementorNumberLargeThenMaximumValue(t *testing.T) {
//	const want = 0
//	t.Run(fmt.Sprintf("%v in %v", MaximumInt, MaximumInt), func(t *testing.T) {
//		var inc Incrementor = NewIncrementor()
//		if err := inc.SetMaximumValue(3); err != nil {
//			t.Fatal(err)
//		}
//
//		// Находим указатель на поле name структуры increment и меняем значение поля
//		// inc.(*increment) - рприведение интерфейса Incrementor к типу increment
//		number := (*int)(unsafe.Pointer(inc.(*increment)))
//		*number = 42
//
//		inc.IncrementNumber()
//		if got := inc.GetNumber(); got != want {
//			t.Errorf("got %v; want %v", got, want)
//		}
//	})
//}
//
//// Тест на установку отрицательно значения через unsafe
//func TestIncrementorNumberNegative(t *testing.T) {
//	const want = -41
//	t.Run(fmt.Sprintf("%v in %v", MaximumInt, MaximumInt), func(t *testing.T) {
//		var inc Incrementor = NewIncrementor()
//		if err := inc.SetMaximumValue(3); err != nil {
//			t.Fatal(err)
//		}
//
//		// Находим указатель на поле name структуры increment и меняем значение поля
//		number := (*int)(unsafe.Pointer(inc.(*increment)))
//		*number = -42
//
//		inc.IncrementNumber()
//		if got := inc.GetNumber(); got != want {
//			t.Errorf("got %v; want %v", got, want)
//		}
//	})
//}
//
////Тест на установку отрицательного значения максимума
//func TestIncrementorMaximumValueNegative(t *testing.T) {
//	const want = 0
//	t.Run(fmt.Sprintf("%v in %v", MaximumInt, MaximumInt), func(t *testing.T) {
//		var inc Incrementor = NewIncrementor()
//		if err := inc.SetMaximumValue(3); err != nil {
//			t.Fatal(err)
//		}
//
//		// Находим указатель на поле name структуры increment и меняем значение поля
//		// Sizeof(42) - для нахождения размера типа int
//		maximumValue := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(inc.(*increment))) + unsafe.Sizeof(42)))
//		*maximumValue = -42
//
//		inc.IncrementNumber()
//		if got := inc.GetNumber(); got != want {
//			t.Errorf("got %v; want %v", got, want)
//		}
//	})
//}


// Бэнчмарк. Просто бэнчмарк. Преждевременная оптимизация кода намеренно не проводилась.
func BenchmarkIncrementor(b *testing.B) {
	inc := NewIncrementor()
	for i := 0; i < b.N; i++ {
		inc.IncrementNumber()
	}
}
