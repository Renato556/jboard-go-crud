package services

import (
	"context"
	"errors"
	"jboard-go-crud/internal/models"
	"jboard-go-crud/internal/repositories"
	"log"
)

type SkillService interface {
	GetAllSkills(ctx context.Context, username string) (models.Skill, error)
	AddSkill(ctx context.Context, skillRequest models.SkillRequest) error
	RemoveSkill(ctx context.Context, skillRequest models.SkillRequest) error
	DeleteUserSkills(ctx context.Context, username string) error
}

type skillService struct {
	skillRepository repositories.SkillRepository
}

func NewSkillService(skillRepository repositories.SkillRepository) SkillService {
	return &skillService{
		skillRepository: skillRepository,
	}
}

func (s *skillService) GetAllSkills(ctx context.Context, username string) (models.Skill, error) {
	log.Printf("Service GetAllSkills called for username: %s", username)

	if username == "" {
		log.Printf("Validation error in GetAllSkills: username is required")
		return models.Skill{}, errors.New("username is required")
	}
	log.Printf("Validation passed for GetAllSkills username: %s", username)

	skill, exists, err := s.skillRepository.FindByUsername(ctx, username)
	if err != nil {
		log.Printf("ERROR: Repository error in GetAllSkills for username %s: %v", username, err)
		return models.Skill{}, err
	}

	if !exists {
		log.Printf("User not found in GetAllSkills for username: %s", username)
		return models.Skill{}, errors.New("user not found")
	}

	log.Printf("Successfully retrieved skills for username: %s", username)
	return skill, nil
}

func (s *skillService) AddSkill(ctx context.Context, skillRequest models.SkillRequest) error {
	log.Printf("Service AddSkill called for username: %s, skill: %s", skillRequest.Username, skillRequest.Skill)

	if skillRequest.Username == "" || skillRequest.Skill == "" {
		log.Printf("Validation error in AddSkill: username and skill are required")
		return errors.New("username and skill are required")
	}
	log.Printf("Validation passed for AddSkill username: %s, skill: %s", skillRequest.Username, skillRequest.Skill)

	// Verifica se o usuário já tem skills cadastradas
	_, exists, err := s.skillRepository.FindByUsername(ctx, skillRequest.Username)
	if err != nil {
		log.Printf("ERROR: Repository error checking user existence in AddSkill for username %s: %v", skillRequest.Username, err)
		return err
	}

	if !exists {
		log.Printf("User does not exist, creating new skills document for username: %s", skillRequest.Username)
		err := s.skillRepository.Create(ctx, skillRequest)
		if err != nil {
			log.Printf("ERROR: Failed to create skills for username %s: %v", skillRequest.Username, err)
		} else {
			log.Printf("Successfully created skills for username: %s", skillRequest.Username)
		}
		return err
	} else {
		log.Printf("User exists, adding skill to existing document for username: %s", skillRequest.Username)
		err := s.skillRepository.AddSkill(ctx, skillRequest)
		if err != nil {
			log.Printf("ERROR: Failed to add skill for username %s: %v", skillRequest.Username, err)
		} else {
			log.Printf("Successfully added skill for username: %s", skillRequest.Username)
		}
		return err
	}
}

func (s *skillService) RemoveSkill(ctx context.Context, skillRequest models.SkillRequest) error {
	log.Printf("Service RemoveSkill called for username: %s, skill: %s", skillRequest.Username, skillRequest.Skill)

	if skillRequest.Username == "" || skillRequest.Skill == "" {
		log.Printf("Validation error in RemoveSkill: username and skill are required")
		return errors.New("username and skill are required")
	}
	log.Printf("Validation passed for RemoveSkill username: %s, skill: %s", skillRequest.Username, skillRequest.Skill)

	// Verifica se o usuário existe
	_, exists, err := s.skillRepository.FindByUsername(ctx, skillRequest.Username)
	if err != nil {
		log.Printf("ERROR: Repository error checking user existence in RemoveSkill for username %s: %v", skillRequest.Username, err)
		return err
	}

	if !exists {
		log.Printf("User not found in RemoveSkill for username: %s", skillRequest.Username)
		return errors.New("user not found")
	}

	log.Printf("User found, proceeding to remove skill for username: %s", skillRequest.Username)
	err = s.skillRepository.RemoveSkill(ctx, skillRequest)
	if err != nil {
		log.Printf("ERROR: Failed to remove skill for username %s: %v", skillRequest.Username, err)
	} else {
		log.Printf("Successfully removed skill for username: %s", skillRequest.Username)
	}
	return err
}

func (s *skillService) DeleteUserSkills(ctx context.Context, username string) error {
	log.Printf("Service DeleteUserSkills called for username: %s", username)

	if username == "" {
		log.Printf("Validation error in DeleteUserSkills: username is required")
		return errors.New("username is required")
	}
	log.Printf("Validation passed for DeleteUserSkills username: %s", username)

	// Verifica se o usuário existe
	_, exists, err := s.skillRepository.FindByUsername(ctx, username)
	if err != nil {
		log.Printf("ERROR: Repository error checking user existence in DeleteUserSkills for username %s: %v", username, err)
		return err
	}

	if !exists {
		log.Printf("User not found in DeleteUserSkills for username: %s", username)
		return errors.New("user not found")
	}

	log.Printf("User found, proceeding to delete all skills for username: %s", username)
	err = s.skillRepository.DeleteByUsername(ctx, username)
	if err != nil {
		log.Printf("ERROR: Failed to delete skills for username %s: %v", username, err)
	} else {
		log.Printf("Successfully deleted all skills for username: %s", username)
	}
	return err
}
