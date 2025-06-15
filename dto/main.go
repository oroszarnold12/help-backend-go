package dto

type Dtoer[T any] interface {
	ToDto() T
}

type ThinDtoer[T any] interface {
	ToThinDto() T
}

func ModelsToDtos[T Dtoer[D], D any](models []T) []D {
	dtos := make([]D, len(models))

	for index := range models {
		dtos[index] = models[index].ToDto()
	}

	return dtos
}

func ModelsToThinDtos[T ThinDtoer[D], D any](models []T) []D {
	dtos := make([]D, len(models))

	for index := range models {
		dtos[index] = models[index].ToThinDto()
	}

	return dtos
}
